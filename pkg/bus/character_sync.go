package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	CharacterSync BusMessageType = "sro.character.sync"
)

type CharacterSyncMessage struct {
	Ids    [][]byte `json:"ids"`
	Batch   int      `json:"batch"`
	traceId string
}

func (m CharacterSyncMessage) GetType() BusMessageType {
	return CharacterSync
}

func (m CharacterSyncMessage) GetId() string {
	return m.traceId
}

func NewCharacterSyncMessage(ctx context.Context, ids [][]byte) CharacterSyncMessage {
	return CharacterSyncMessage{
		Ids:    ids,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
