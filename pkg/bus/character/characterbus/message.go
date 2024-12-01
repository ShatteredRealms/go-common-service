package characterbus

import "github.com/ShatteredRealms/go-common-service/pkg/bus"

type Message struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Id is the unique identifier of the owner of the character
	OwnerId string `json:"ownerId"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.character")
}

func (m Message) GetId() string {
	return m.Id
}
