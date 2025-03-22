package channelbus

import (
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Message struct {
	// Id is the unique identifier of the channel
	Id uuid.UUID `json:"id"`

	// DimensionId is the unique identifier of the dimension that this channel is assigned to.
	// If this is all zeros (empty uuid), then this channel is not assigned to any dimension and is a global channel.
	DimensionId uuid.UUID `json:"dimensionId"`

	// Deleted is a flag indicating if the channel has been deleted
	Deleted bool `json:"deleted"`
}

type BusReader bus.MessageBusReader[Message]
type BusWriter bus.MessageBusWriter[Message]

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.chat.channel")
}

func (m Message) GetId() string {
	return m.Id.String()
}

func (m Message) WasDeleted() bool {
	return m.Deleted
}
