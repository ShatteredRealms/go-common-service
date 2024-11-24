package bus

import "context"

type BusMessageType string

type BusMessage[T any] interface {
	GetType() BusMessageType
	GetId() string
	GetData() T
}

type MessageBus[T any] interface {
	ReceiveMessages(context.Context, chan T) error
	Publish(context.Context, BusMessage[T]) error
	Close(context.Context) error
}

