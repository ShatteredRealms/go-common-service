package bus_test

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

// TestBusMessage is a test message for the bus. It full control
// of the type, id, and deleted status allowing unique topics (types)
// per test call.
type TestBusMessage struct {
	Type    bus.BusMessageType
	Id      string
	Deleted bool
}

func (t TestBusMessage) GetType() bus.BusMessageType {
	if t.Type == "" {
		return bus.BusMessageType("test")
	}
	return t.Type
}

func (t TestBusMessage) GetId() string {
	return t.Id
}
func (t TestBusMessage) WasDeleted() bool {
	return t.Deleted
}

type TestingBus struct {
	ErrOnClose            error
	ErrOnFetch            error
	ErrOnProcessFailed    error
	ErrOnProcessSucceeded error
	ErrOnProcessSkipped   error
	ErrOnPublish          error
	Queue                 chan *TestBusMessage
	CurrentMessage        *TestBusMessage
}

// Publish implements bus.MessageBusWriter.
func (t *TestingBus) Publish(ctx context.Context, msg TestBusMessage) error {
	if t.ErrOnPublish != nil {
		return t.ErrOnPublish
	}
	if len(t.Queue) == cap(t.Queue) {
		return fmt.Errorf("queue is full")
	}
	t.Queue <- &msg
	return nil
}

// Close implements bus.MessageBusReader.
func (t *TestingBus) Close() error {
	return t.ErrOnClose
}

// FetchMessage implements bus.MessageBusReader.
func (t *TestingBus) FetchMessage(context.Context) (*TestBusMessage, error) {
	if t.ErrOnFetch != nil {
		return nil, t.ErrOnFetch
	}

	if t.CurrentMessage == nil {
		t.CurrentMessage = <-t.Queue
	}

	return t.CurrentMessage, nil
}

// ProcessFailed implements bus.MessageBusReader.
func (t *TestingBus) ProcessFailed() error {
	return t.ErrOnProcessFailed
}

// ProcessSucceeded implements bus.MessageBusReader.
func (t *TestingBus) ProcessSucceeded(context.Context) error {
	t.CurrentMessage = nil
	return t.ErrOnProcessSucceeded
}

// ProcessSucceeded implements bus.MessageBusReader.
func (t *TestingBus) ProcessSkipped(context.Context) error {
	t.CurrentMessage = nil
	return t.ErrOnProcessSkipped
}

func (t *TestingBus) Reset(context.Context) error {
	t.CurrentMessage = nil
	return nil
}

func (t *TestingBus) GetMessageType() bus.BusMessageType {
	return bus.BusMessageType("test")
}

func (t *TestingBus) GetGroup() string {
	return "test"
}
