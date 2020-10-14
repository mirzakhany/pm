package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/cycles"
)

type Cycle struct {
	tableName   struct{} `pg:"cycles,alias:i"` //nolint
	ID          uint64
	UUID        string `pg:"default:gen_random_uuid()"`
	Title       string
	Description string
	Active      bool
	StartAt     time.Time
	EndAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (cm Cycle) ToProto(secure bool) *cycles.Cycle {
	c, _ := ptypes.TimestampProto(cm.CreatedAt)
	u, _ := ptypes.TimestampProto(cm.UpdatedAt)

	s, _ := ptypes.TimestampProto(cm.StartAt)
	e, _ := ptypes.TimestampProto(cm.EndAt)

	cycle := &cycles.Cycle{
		Id:          cm.ID,
		Uuid:        cm.UUID,
		Title:       cm.Title,
		Description: cm.Description,
		Active:      cm.Active,
		StartAt:     s,
		EndAt:       e,
		CreatedAt:   c,
		UpdatedAt:   u,
	}
	return cycle
}

func CycleToProtoList(cml []Cycle, secure bool) []*cycles.Cycle {
	var c []*cycles.Cycle
	for _, i := range cml {
		c = append(c, i.ToProto(secure))
	}
	return c
}

func CycleFromProto(cycle *cycles.Cycle) Cycle {
	c, _ := ptypes.Timestamp(cycle.CreatedAt)
	u, _ := ptypes.Timestamp(cycle.UpdatedAt)

	s, _ := ptypes.Timestamp(cycle.StartAt)
	e, _ := ptypes.Timestamp(cycle.EndAt)

	return Cycle{
		ID:          cycle.Id,
		UUID:        cycle.Uuid,
		Title:       cycle.Title,
		Description: cycle.Description,
		Active:      cycle.Active,
		StartAt:     s,
		EndAt:       e,
		CreatedAt:   c,
		UpdatedAt:   u,
	}
}
