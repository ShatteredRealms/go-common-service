package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	CharacterDeleted BusMessageType = "sro.character.deleted"
)

type CharacterDeletedMessage struct {
	ID      string `json:"id"`
	traceId string
}

func (m CharacterDeletedMessage) GetType() BusMessageType {
	return CharacterDeleted
}

func (m CharacterDeletedMessage) GetId() string {
	return m.traceId
}

func NewCharacterDeletedMessage(ctx context.Context, id string) CharacterDeletedMessage {
	return CharacterDeletedMessage{
		ID:      id,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
