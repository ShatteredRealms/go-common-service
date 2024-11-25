package bus

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	DimensionCreated BusMessageType = "sro.dimension.created"
)

type DimensionCreatedMessage struct {
	ID      string `json:"id"`
	traceId string
}

func (m DimensionCreatedMessage) GetType() BusMessageType {
	return DimensionCreated
}

func (m DimensionCreatedMessage) GetId() string {
	return m.traceId
}

func NewDimensionCreatedMessage(ctx context.Context, id string) DimensionCreatedMessage {
	return DimensionCreatedMessage{
		ID:      id,
		traceId: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
