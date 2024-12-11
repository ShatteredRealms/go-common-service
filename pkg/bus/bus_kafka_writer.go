package bus

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

// MessageBusWriter is for writing message asynchronously to the message bus.
type kafkaBusWriter[T BusMessage[any]] struct {
	*kafkaBus[T]
	Writer *kafka.Writer
}

// Publish implements MessageBus.
func (k *kafkaBusWriter[T]) Publish(ctx context.Context, msg T) error {
	ctx, span := k.tracer.Start(ctx, "Publish")
	defer span.End()
	span.SetAttributes(BusMessageAttributes(msg)...)
	k.setupWriter()

	val, err := k.encodeMessage(msg)
	if err != nil {
		return err
	}

	ctx, innerSpan := k.tracer.Start(ctx, "Publish.SendMessages")
	defer innerSpan.End()
	k.wg.Add(1)
	defer k.wg.Done()
	err = k.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.GetId()),
		Value: val,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendingMessage, err)
	}

	return nil
}

func (k *kafkaBusWriter[T]) PublishMany(ctx context.Context, msgs []T) error {
	k.wg.Add(1)
	defer k.wg.Done()

	if len(msgs) == 0 {
		return nil
	}

	ctx, span := k.tracer.Start(ctx, "PublishMany")
	span.SetAttributes(BusMessageAttributes(msgs[0])[0])
	defer span.End()

	k.setupWriter()

	messages := make([]kafka.Message, len(msgs))
	messageMu := sync.Mutex{}
	var errs error
	errsMu := sync.Mutex{}
	wg := sync.WaitGroup{}

	ctx, innerSpan := k.tracer.Start(ctx, "PublishMany.CreateMessages")
	for chunk := range slices.Chunk(msgs, runtime.NumCPU()) {
		wg.Add(1)
		go func(chunk []T) {
			defer wg.Done()
			for _, msg := range chunk {
				val, err := k.encodeMessage(msg)
				if err != nil {
					errsMu.Lock()
					errs = errors.Join(errs, fmt.Errorf("%w: %w", ErrEncodingMessage, err))
					errsMu.Unlock()
					return
				}

				messageMu.Lock()
				messages = append(messages, kafka.Message{
					Key:   []byte(msg.GetId()),
					Value: val,
				})
				messageMu.Unlock()
			}
		}(chunk)
	}

	wg.Wait()
	innerSpan.End()
	if errs != nil {
		return errs
	}

	ctx, innerSpan = k.tracer.Start(ctx, "PublishMany.SendMessages")
	defer innerSpan.End()
	err := k.Writer.WriteMessages(ctx, messages...)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendingMessage, err)
	}

	return nil
}

func (k *kafkaBusWriter[T]) GetMessageType() BusMessageType {
	return BusMessageType(k.topic)
}

func (k *kafkaBusWriter[T]) Close() error {
	k.wg.Wait()

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.Writer != nil {
		err := k.Writer.Close()
		k.Writer = nil
		return err
	}

	return nil
}

func (k *kafkaBusWriter[T]) setupWriter() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.Writer == nil {
		k.Writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  k.brokers.Addresses(),
			Topic:    k.topic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
			Logger:   kafka.LoggerFunc(log.Logger.Tracef),
		})
		k.Writer.AllowAutoTopicCreation = true
	}
}

func (k *kafkaBusWriter[T]) encodeMessage(msg T) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}
	return buf.Bytes(), nil
}

func NewKafkaMessageBusWriter[T BusMessage[any]](brokers config.ServerAddresses, msg T) MessageBusWriter[T] {
	return &kafkaBusWriter[T]{
		kafkaBus: &kafkaBus[T]{
			brokers: brokers,
			topic:   string(msg.GetType()),
			tracer:  otel.Tracer(fmt.Sprintf("sro.bus.kafka.reader.%s", msg.GetType())),
		},
	}
}
