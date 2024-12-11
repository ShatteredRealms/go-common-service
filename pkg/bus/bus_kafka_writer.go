package bus

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
)

// MessageBusWriter is for writing message asynchronously to the message bus.
type kafkaBusWriter[T BusMessage[any]] struct {
	*kafkaBus[T]
	Writer *kafka.Writer
}

// Publish implements MessageBus.
func (k *kafkaBusWriter[T]) Publish(ctx context.Context, msg T) error {
	k.mu.Lock()
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
	k.mu.Unlock()

	k.wg.Add(1)
	defer k.wg.Done()
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}

	err = k.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.GetId()),
		Value: buf.Bytes(),
	})
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
	if k.Writer != nil {
		err := k.Writer.Close()
		k.Writer = nil
		return err
	}
	return nil
}

func NewKafkaMessageBusWriter[T BusMessage[any]](brokers config.ServerAddresses, msg T) MessageBusWriter[T] {
	return &kafkaBusWriter[T]{
		kafkaBus: &kafkaBus[T]{
			brokers: brokers,
			topic:   string(msg.GetType()),
		},
	}
}
