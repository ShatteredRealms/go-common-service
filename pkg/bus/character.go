package bus

const (
	Character BusMessageType = "sro.character"
)

type CharacterMessage struct {
	// Id is the unique identifier of the character
	Id string `json:"id"`

	// Id is the unique identifier of the owner of the character
	OwnerId string `json:"ownerId"`

	// Deleted is a flag indicating if the character has been deleted
	Deleted bool `json:"deleted"`
}

func (m CharacterMessage) GetType() BusMessageType {
	return Character
}

func (m CharacterMessage) GetId() string {
	return m.Id
}
