package bus

import "context"

type BusMessageRepository[T BusModelMessage[any]] interface {
	Save(ctx context.Context, data T) error
	Delete(ctx context.Context, id string) error
}
