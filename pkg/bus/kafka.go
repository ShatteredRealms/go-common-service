package bus

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
)

var (
	ErrSerializeMessage = errors.New("unable to serialize message")
	ErrSendingMessage   = errors.New("unable to send message on bus")
)

type kafkaBus[T any] struct {
	Reader  *kafka.Reader
	Writer  *kafka.Writer
	brokers config.ServerAddresses
	groupId string
	topic   string
}

// Publish implements MessageBus.
func (k *kafkaBus[T]) Publish(ctx context.Context, msg BusMessage[T]) error {
	k.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.brokers.Addresses(),
		Topic:    k.topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
		Logger:   kafka.LoggerFunc(log.Logger.Debugf),
	})
	k.Writer.AllowAutoTopicCreation = true
	defer k.Writer.Close()

	key, err := toByteArray(msg.GetId())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}

	data, err := toByteArray(msg.GetData())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSerializeMessage, err)
	}

	err = k.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendingMessage, err)
	}

	return nil
}

// ReceiveMessages implements MessageBus.
func (k *kafkaBus[T]) ReceiveMessages(ctx context.Context, channel chan T) error {
	k.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:               k.brokers.Addresses(),
		Topic:                 k.topic,
		GroupID:        k.groupId,
		MinBytes:       1,
		MaxBytes:       10e3,
		CommitInterval: 1 * time.Second,
		Logger:         kafka.LoggerFunc(log.Logger.Debugf),
	})
	defer k.Reader.Close()

	for ctx.Err() == nil {
		msg, err := k.Reader.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("%w: %w", ErrSerializeMessage, err)
		} else {
			var data T
			dec := gob.NewDecoder(bytes.NewReader(msg.Value))
			if err := dec.Decode(&data); err != nil {
				return fmt.Errorf("%w: %w", ErrSerializeMessage, err)
			} else {
				channel <- data
			}
		}
	}
	return nil
}

func (k *kafkaBus[T]) Close(ctx context.Context) error {
	var errs error
	if k.Reader != nil {
		errors.Join(errs, k.Reader.Close())
	}
	if k.Writer != nil {
		errors.Join(errs, k.Writer.Close())
	}
	return errs
}

func NewKafkaMessageBus[T any](brokers config.ServerAddresses, groupId string, msg BusMessage[T]) MessageBus[T] {
	return &kafkaBus[T]{
		brokers: brokers,
		topic:   string(msg.GetType()),
		groupId: groupId,
	}
}

func toByteArray[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
