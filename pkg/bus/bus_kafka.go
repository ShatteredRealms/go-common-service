package bus

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
)

var (
	ErrSerializeMessage = errors.New("unable to serialize message")
	ErrSendingMessage   = errors.New("unable to send message on bus")
)

type kafkaBus[T BusMessage[any]] struct {
	brokers config.ServerAddresses
	topic   string
}

// MessageBusReader is for reading messages from the message bus synchronously.
type kafkaBusReader[T BusMessage[any]] struct {
	*kafkaBus[T]
	groupId        string
	Reader         *kafka.Reader
	currentMessage *kafka.Message
}

// MessageBusWriter is for writing message asynchronously to the message bus.
type kafkaBusWriter[T BusMessage[any]] struct {
	*kafkaBus[T]
	Writer *kafka.Writer

	mu sync.Mutex
	wg sync.WaitGroup
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

// ReceiveMessages implements MessageBus.
func (k *kafkaBusReader[T]) FetchMessage(ctx context.Context) (*T, error) {
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

func (k *kafkaBusReader[T]) Close() error {
	if k.Reader != nil {
		err := k.Reader.Close()
		k.Reader = nil
		return err
	}
	return nil
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

func NewKafkaMessageBusReader[T BusMessage[any]](brokers config.ServerAddresses, groupId string, msg T) MessageBusReader[T] {
	return &kafkaBusReader[T]{
		kafkaBus: &kafkaBus[T]{
			brokers: brokers,
			topic:   string(msg.GetType()),
		},
		groupId: groupId,
	}
}

func NewKafkaMessageBusWriter[T BusMessage[any]](brokers config.ServerAddresses, msg T) MessageBusWriter[T] {
	return &kafkaBusWriter[T]{
		kafkaBus: &kafkaBus[T]{
			brokers: brokers,
			topic:   string(msg.GetType()),
		},
	}
}
