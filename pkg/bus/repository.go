package bus

import "context"

type BusMessageRepository interface {
	Save(ctx context.Context, data any) error
	Delete(ctx context.Context, id string) error
}
