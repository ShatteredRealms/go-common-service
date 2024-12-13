package dimensionbus

import (
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Message struct {
	// Id is the unique identifier of the dimension
	Id uuid.UUID `json:"id"`

	// Deleted is a flag indicating if the dimension has been deleted
	Deleted bool `json:"deleted"`
}

type BusReader bus.MessageBusReader[Message]
type BusWriter bus.MessageBusWriter[Message]

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.gameserver.dimension")
}

func (m Message) GetId() string {
	return m.Id.String()
}

func (m Message) WasDeleted() bool {
	return m.Deleted
}
