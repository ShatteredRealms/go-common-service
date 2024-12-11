package mapbus

import (
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Message struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.gameserver.map")
}

func (m Message) GetId() string {
	return m.Id
}

func (m Message) WasDeleted() bool {
	return m.Deleted
}
