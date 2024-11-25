package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	CharacterCreated BusMessageType = "sro.character.created"
)

type CharacterCreatedMessage struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	traceId string
}

func (m CharacterCreatedMessage) GetType() BusMessageType {
	return CharacterCreated
}

func (m CharacterCreatedMessage) GetId() string {
	return m.traceId
}

func NewCharacterCreatedMessage(ctx context.Context, id, name string) CharacterCreatedMessage {
	return CharacterCreatedMessage{
		ID:      id,
		Name:    name,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
