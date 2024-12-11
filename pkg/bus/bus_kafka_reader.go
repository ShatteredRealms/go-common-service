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
	if k.isResetting {
		log.Logger.WithContext(ctx).Info("Waiting for reset to finish")
		timer := time.NewTimer(30 * time.Second)
		select {
		case <-k.resetFinished:
		case <-timer.C:
			if k.isResetting {
				return nil, errors.New("resetting took too long")
			}
		}
	}

	if k.currentMessage != nil {
		return nil, errors.New("message already fetched")
	}
	k.currentMessage = new(kafka.Message)

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
	var err error
	*k.currentMessage, err = k.Reader.FetchMessage(ctx)
	if k.isResetting {
		log.Logger.WithContext(ctx).Info("Reset started, skipping message")
		return k.FetchMessage(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}

	var data T
	dec := gob.NewDecoder(bytes.NewReader(k.currentMessage.Value))
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}

	return &data, nil
}

func (k *kafkaBusReader[T]) Reset(ctx context.Context) error {
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
		k.Reader = nil
		k.currentMessage = nil
	}

	// Connect to the kafka cluster
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Connecting to Kafka")
	conn, err := kafka.DialContext(ctx, "tcp", k.brokers.Addresses()[0])
	if err != nil {
		return fmt.Errorf("unable to dial kafka: %w", err)
	}
	defer conn.Close()

	// Read the partitions for the topic
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Reading partitions")
	partitions, err := conn.ReadPartitions(k.topic)
	if err != nil {
		return fmt.Errorf("unable to read partitions: %w", err)
	}

	wg := sync.WaitGroup{}
	errMu := sync.Mutex{}
	offsetMu := sync.Mutex{}

	var outErrors error
	offsets := make(map[int]int64, len(partitions))

	// Get the first offset for each partition
	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Infof("Getting offsets for %d partitions", len(partitions))
	for _, partition := range partitions {
		wg.Add(1)
		go func() {
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
		}()
	}

	// Wait for all the offsets to be read
	wg.Wait()
	if outErrors != nil {
		return outErrors
	}

	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Infof("Getting consumer group %s", k.groupId)
	// Set offset to the beginning for this consumer group and topic
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      k.groupId,
		Brokers: k.brokers.Addresses(),
		Topics:  []string{k.topic},
	})
	if err != nil {
		return fmt.Errorf("unable to create consumer group: %w", err)
	}
	defer group.Close()

	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Getting next generation")
	generation, err := group.Next(ctx)
	if err != nil {
		return fmt.Errorf("unable to get next generation: %w", err)
	}

	log.Logger.WithContext(ctx).WithField("func", "Bus Reset").Info("Commiting offsets")
	err = generation.CommitOffsets(map[string]map[int]int64{k.topic: offsets})
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
		},
		groupId: groupId,
	}
}
