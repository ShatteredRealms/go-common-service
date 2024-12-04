package bus_test

import (
	"context"
)

type TestBusMessageRepository struct {
	ErrOnSave   error
	ErrOnDelete error
}

// Delete implements bus.BusMessageRepository.
func (t *TestBusMessageRepository) Delete(ctx context.Context, id string) error {
	return t.ErrOnDelete
}

// Save implements bus.BusMessageRepository.
func (t *TestBusMessageRepository) Save(ctx context.Context, data TestBusMessage) error {
	return t.ErrOnSave
}

