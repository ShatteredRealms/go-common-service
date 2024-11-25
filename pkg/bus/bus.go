package bus

import "context"

type BusMessageType string

type BusMessage[T any] interface {
	GetType() BusMessageType
	GetId() string
}

type MessageBusReader[T BusMessage[any]] interface {
	FetchMessage(context.Context) (*T, error)
	ProcessSucceeded(context.Context) error
	ProcessFailed(context.Context) error
	Close(context.Context) error
}

type MessageBusWriter[T BusMessage[any]] interface {
	Publish(context.Context, T) error
	Close(context.Context) error
}
