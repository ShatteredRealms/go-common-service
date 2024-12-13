package bus

import (
	"context"

	"github.com/google/uuid"
)

type BusMessageRepository[T BusModelMessage[any]] interface {
	Save(ctx context.Context, data T) error
	Delete(ctx context.Context, id *uuid.UUID) error
}
