package bus

import "context"

type BusMessageType string

type BusMessage[T any] interface {
	GetType() BusMessageType
	GetId() string
}

type BusModelMessage[T any] interface {
	BusMessage[T]
	WasDeleted() bool
}

type MessageBusReader[T BusMessage[any]] interface {
	FetchMessage(context.Context) (*T, error)
	ProcessSucceeded(context.Context) error
	ProcessFailed() error
	Close() error
}

type MessageBusWriter[T BusMessage[any]] interface {
	Publish(context.Context, T) error
	Close() error
}
