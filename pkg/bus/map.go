package bus

const (
	Map BusMessageType = "sro.gameserver.map"
)

type MapMessage struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

func (m MapMessage) GetType() BusMessageType {
	return Map
}

func (m MapMessage) GetId() string {
	return m.Id
}
