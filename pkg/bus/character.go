package bus

const (
	CharacterCreated BusMessageType = "sro.character.created"
)

type CharacterMessage struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

func (m CharacterMessage) GetType() BusMessageType {
	return CharacterCreated
}

func (m CharacterMessage) GetId() string {
	return m.Id
}
