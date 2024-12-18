package game

import (
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/google/uuid"
)

type Location struct {
	WorldId uuid.UUID `db:"world_id" json:"worldId" mapstructure:"world_id"`
	X       float32   `db:"x" json:"x" mapstructure:"x"`
	Y       float32   `db:"y" json:"y" mapstructure:"y"`
	Z       float32   `db:"z" json:"z" mapstructure:"z"`
	Roll    float32   `db:"roll" json:"roll" mapstructure:"roll"`
	Pitch   float32   `db:"pitch" json:"pitch" mapstructure:"pitch"`
	Yaw     float32   `db:"yaw" json:"yaw" mapstructure:"yaw"`
}

func (l Location) ToPb() *pb.Location {
	return &pb.Location{
		World: l.WorldId.String(),
		X:     l.X,
		Y:     l.Y,
		Z:     l.Z,
		Roll:  l.Roll,
		Pitch: l.Pitch,
		Yaw:   l.Yaw,
	}
}

func LocationFromPb(location *pb.Location) (*Location, error) {
	worldId, err := uuid.Parse(location.World)
	return &Location{
		WorldId: worldId,
		X:       location.X,
		Y:       location.Y,
		Z:       location.Z,
		Roll:    location.Roll,
		Pitch:   location.Pitch,
		Yaw:     location.Yaw,
	}, err
}
