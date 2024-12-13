package characterbus

import (
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Message struct {
	// Id is the unique identifier of the character
	Id uuid.UUID `json:"id"`

	// OwnerId is the unique identifier of the owner of the character
	OwnerId uuid.UUID `json:"ownerId"`

	// DimensionId is the unique identifier of the dimension the character is in
	DimensionId uuid.UUID `json:"dimensionId"`

	// MapId is the unique identifier of the map the character is in
	MapId uuid.UUID `json:"mapId"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

type BusReader bus.MessageBusReader[Message]
type BusWriter bus.MessageBusWriter[Message]

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.character")
}

func (m Message) GetId() string {
	return m.Id.String()
}

func (m Message) WasDeleted() bool {
	return m.Deleted
}
