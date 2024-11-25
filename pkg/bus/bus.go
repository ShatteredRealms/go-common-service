package bus

import "context"

type BusMessageType string

type BusMessage[T any] interface {
	GetType() BusMessageType
	GetId() string
}

type MessageBus[T BusMessage[any]] interface {
	ReceiveMessages(context.Context, chan T) error
	Publish(context.Context, T) error
	Close(context.Context) error
}

