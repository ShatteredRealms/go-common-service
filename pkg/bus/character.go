package bus

const (
	Character BusMessageType = "sro.character"
)

type CharacterMessage struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

func (m CharacterMessage) GetType() BusMessageType {
	return Character
}

func (m CharacterMessage) GetId() string {
	return m.Id
}
