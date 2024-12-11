package characterbus

import "github.com/ShatteredRealms/go-common-service/pkg/bus"

type Message struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// OwnerId is the unique identifier of the owner of the character
	OwnerId string `json:"ownerId"`

	// DimensionId is the unique identifier of the dimension the character is in
	DimensionId string `json:"dimensionId"`

	// MapId is the unique identifier of the map the character is in
	MapId string `json:"mapId"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

type BusReader bus.MessageBusReader[Message]
type BusWriter bus.MessageBusWriter[Message]

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.character")
}

func (m Message) GetId() string {
	return m.Id
}

func (m Message) WasDeleted() bool {
	return m.Deleted
}
