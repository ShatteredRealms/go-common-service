package dimensionbus

import "github.com/ShatteredRealms/go-common-service/pkg/bus"

type Message struct {
	// Id is the unique identifier of the dimension
	Id string `json:"id"`

	// Deleted is a flag indicating if the dimension has been deleted
	Deleted bool `json:"deleted"`
}

func (m Message) GetType() bus.BusMessageType {
	return bus.BusMessageType("sro.gameserver.dimension")
}

func (m Message) GetId() string {
	return m.Id
}
