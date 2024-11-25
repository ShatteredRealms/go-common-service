package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	DimensionDeleted BusMessageType = "sro.dimension.deleted"
)

type DimensionDeletedMessage struct {
	ID      string `json:"id"`
	traceId string
}

func (m DimensionDeletedMessage) GetType() BusMessageType {
	return CharacterSync
}

func (m DimensionDeletedMessage) GetId() string {
	return m.traceId
}

func NewDimensionDeletedMessage(ctx context.Context, id string) DimensionDeletedMessage {
	return DimensionDeletedMessage{
		ID:      id,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
