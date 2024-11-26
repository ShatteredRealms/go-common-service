package bus

const (
	Dimension BusMessageType = "sro.gameserver.dimension"
)

type DimensionMessage struct {
	// Id is the unique identifier of the dimension
	Id string `json:"id"`

	// Deleted is a flag indicating if the dimension has been deleted
	Deleted bool `json:"deleted"`
}

func (m DimensionMessage) GetType() BusMessageType {
	return Dimension
}

func (m DimensionMessage) GetId() string {
	return m.Id
}
