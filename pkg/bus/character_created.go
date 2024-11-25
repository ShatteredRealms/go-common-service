package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	CharacterCreated BusMessageType = "sro.character.created"
)

type CharacterCreatedMessage struct {
	Id      string `json:"id"`
	traceId string
}

func (m CharacterCreatedMessage) GetType() BusMessageType {
	return CharacterCreated
}

func (m CharacterCreatedMessage) GetId() string {
	return m.Id
}

func NewCharacterCreatedMessage(ctx context.Context, id string) CharacterCreatedMessage {
	return CharacterCreatedMessage{
		Id:      id,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
