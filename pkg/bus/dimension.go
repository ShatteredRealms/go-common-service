package bus

const (
	DimensionCreated BusMessageType = "sro.dimension.created"
)

type DimensionMessage struct {
	// Id is the unique identifier of the dimension
	Id string `json:"id"`

	// Deleted is a flag indicating if the dimension has been deleted
	Deleted bool `json:"deleted"`
}

func (m DimensionMessage) GetType() BusMessageType {
	return DimensionCreated
}

func (m DimensionMessage) GetId() string {
	return m.Id
}
