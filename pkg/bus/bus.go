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

type Resettable interface {
	Reset(context.Context) error
}

type MessageBusReader[T BusMessage[any]] interface {
	Resettable
	GetMessageType() BusMessageType
	GetGroup() string
	FetchMessage(context.Context) (*T, error)
	ProcessSucceeded(context.Context) error
	ProcessFailed() error
	Close() error
}

type MessageBusWriter[T BusMessage[any]] interface {
	GetMessageType() BusMessageType
	Publish(context.Context, T) error
	PublishMany(context.Context, []T) error
	Close() error
}
