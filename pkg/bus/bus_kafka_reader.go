package bus

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

// MessageBusReader is for reading messages from the message bus synchronously.
type kafkaBusReader[T BusMessage[any]] struct {
	*kafkaBus[T]
	groupId        string
	Reader         *kafka.Reader
	currentMessage *kafka.Message

	isResetting   bool
	resetFinished chan struct{}
}

// ReceiveMessages implements MessageBus.
func (k *kafkaBusReader[T]) FetchMessage(ctx context.Context) (*T, error) {
	ctx, span := k.tracer.Start(ctx, "fetch")
	defer span.End()
	if k.isResetting {
		ctx, innerSpan := k.tracer.Start(ctx, "fetch")
		log.Logger.WithContext(ctx).Info("Waiting for reset to finish")
		timer := time.NewTimer(30 * time.Second)
		select {
		case <-k.resetFinished:
		case <-timer.C:
			if k.isResetting {
				innerSpan.End()
				return nil, errors.New("resetting took too long")
			}
		}
		innerSpan.End()
	}

	if k.currentMessage != nil {
		return nil, errors.New("message already fetched")
	}
	k.currentMessage = new(kafka.Message)

	k.mu.Lock()
	if k.Reader == nil {
		k.Reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  k.brokers.Addresses(),
			Topic:    k.topic,
			GroupID:  k.groupId,
			MinBytes: 1,
			MaxBytes: 10e3,
			Logger:   kafka.LoggerFunc(log.Logger.Tracef),
		})
	}

	ctx, innerSpan := k.tracer.Start(ctx, "fetch.message")
	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		*k.currentMessage, err = k.Reader.FetchMessage(ctx)
		wg.Done()
	}()
	k.mu.Unlock()
	wg.Wait()
	innerSpan.End()

	if k.isResetting {
		log.Logger.WithContext(ctx).Info("Reset started, skipping message")
		return k.FetchMessage(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFetchMessage, err)
	}

	ctx, innerSpan = k.tracer.Start(ctx, "fetch.decode")
	defer innerSpan.End()
	var data T
	dec := gob.NewDecoder(bytes.NewReader(k.currentMessage.Value))
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodingMessage, err)
	}

	return &data, nil
}

func (k *kafkaBusReader[T]) Reset(ctx context.Context) error {
	ctx, span := k.tracer.Start(ctx, "reset")
	defer span.End()
	// Prevent multiple resets from happening at the same time
	k.mu.Lock()
	defer k.mu.Unlock()

	k.isResetting = true
	defer func() {
		k.isResetting = false
	}()

	// Close the reader if it is open and cancel the current message
	if k.Reader != nil {
		k.Reader.Close()
		defer func() {
			k.Reader = nil
			k.currentMessage = nil
		}()
	}

	// Connect to the kafka cluster
	ctx, innerSpan := k.tracer.Start(ctx, "reset.connect")
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Connecting to Kafka")
	conn, err := kafka.DialContext(ctx, "tcp", k.brokers.Addresses()[0])
	if err != nil {
		return fmt.Errorf("unable to dial kafka: %w", err)
	}
	defer conn.Close()
	innerSpan.End()

	// Read the partitions for the topic
	ctx, innerSpan = k.tracer.Start(ctx, "reset.partitions.get")
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Reading partitions")
	partitions, err := conn.ReadPartitions(k.topic)
	innerSpan.End()
	if err != nil {
		return fmt.Errorf("unable to read partitions: %w", err)
	}

	ctx, innerSpan = k.tracer.Start(ctx, "reset.partitions.offsets")
	wg := sync.WaitGroup{}
	errMu := sync.Mutex{}
	offsetMu := sync.Mutex{}

	var outErrors error
	offsets := make(map[int]int64, len(partitions))

	// Get the first offset for each partition
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Infof("Getting offsets for %d partitions", len(partitions))
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition kafka.Partition) {
			defer wg.Done()
			leaderAddr := fmt.Sprintf("%s:%d", partition.Leader.Host, partition.Leader.Port)
			c, err := kafka.DialLeader(ctx, "tcp", leaderAddr, k.topic, partition.ID)
			if err != nil {
				errMu.Lock()
				outErrors = errors.Join(outErrors, fmt.Errorf("unable to dial leader: %w", err))
				errMu.Unlock()
				return
			}
			defer c.Close()

			offset, err := c.ReadFirstOffset()
			if err != nil {
				errMu.Lock()
				outErrors = errors.Join(outErrors, fmt.Errorf("unable to read first offset: %w", err))
				errMu.Unlock()
				return
			}

			offsetMu.Lock()
			offsets[partition.ID] = offset
			offsetMu.Unlock()
		}(partition)
	}

	// Wait for all the offsets to be read
	wg.Wait()
	innerSpan.End()
	if outErrors != nil {
		return outErrors
	}

	ctx, innerSpan = k.tracer.Start(ctx, "reset.consumergroup.get")
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Infof("Getting consumer group %s", k.groupId)
	// Set offset to the beginning for this consumer group and topic
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      k.groupId,
		Brokers: k.brokers.Addresses(),
		Topics:  []string{k.topic},
	})
	innerSpan.End()
	if err != nil {
		return fmt.Errorf("unable to create consumer group: %w", err)
	}
	defer group.Close()

	ctx, innerSpan = k.tracer.Start(ctx, "reset.consumergroup.generation.next")
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Getting next generation")
	generation, err := group.Next(ctx)
	innerSpan.End()
	if err != nil {
		return fmt.Errorf("unable to get next generation: %w", err)
	}

	ctx, innerSpan = k.tracer.Start(ctx, "reset.consumergroup.offsets.commit")
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Commiting offsets")
	err = generation.CommitOffsets(map[string]map[int]int64{k.topic: offsets})
	innerSpan.End()
	if err != nil {
		return fmt.Errorf("unable to commit offsets: %w", err)
	}

	return nil
}

func (k *kafkaBusReader[T]) GetMessageType() BusMessageType {
	return BusMessageType(k.topic)
}

func (k *kafkaBusReader[T]) ProcessSucceeded(ctx context.Context) error {
	if k.currentMessage == nil {
		return errors.New("message not fetched")
	}
	if k.Reader != nil {
		err := k.Reader.CommitMessages(ctx, *k.currentMessage)
		k.currentMessage = nil
		return err
	}
	return errors.New("reader not initialized")
}

func (k *kafkaBusReader[T]) ProcessSkipped(ctx context.Context) error {
	if k.currentMessage == nil {
		return errors.New("message not fetched")
	}
	if k.Reader != nil {
		err := k.Reader.CommitMessages(ctx, *k.currentMessage)
		k.currentMessage = nil
		return err
	}
	return errors.New("reader not initialized")
}

func (k *kafkaBusReader[T]) ProcessFailed() error {
	if k.Reader == nil {
		return errors.New("reader not initialized")
	}
	if k.currentMessage == nil {
		return errors.New("message not fetched")
	}
	err := k.Close()
	k.currentMessage = nil
	return err
}

func (k *kafkaBusReader[T]) GetGroup() string {
	return k.groupId
}

func (k *kafkaBusReader[T]) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.Reader != nil {
		err := k.Reader.Close()
		k.Reader = nil
		return err
	}
	return nil
}

func NewKafkaMessageBusReader[T BusMessage[any]](brokers config.ServerAddresses, groupId string, msg T) MessageBusReader[T] {
	return &kafkaBusReader[T]{
		kafkaBus: &kafkaBus[T]{
			brokers: brokers,
			topic:   string(msg.GetType()),
			tracer:  otel.Tracer(fmt.Sprintf("sro.bus.kafka.reader.%s", msg.GetType())),
		},
		groupId: groupId,
	}
}
