package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/roles"
)

type Role struct {
	tableName struct{} `pg:"roles,alias:r"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (rm Role) ToProto() *roles.Role {
	c, _ := ptypes.TimestampProto(rm.CreatedAt)
	u, _ := ptypes.TimestampProto(rm.UpdatedAt)

	role := &roles.Role{
		Id:        rm.ID,
		Uuid:      rm.UUID,
		Title:     rm.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return role
}

func RoleToProtoList(rml []Role) []*roles.Role {
	var r []*roles.Role
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func RoleFromProto(role *roles.Role) Role {
	c, _ := ptypes.Timestamp(role.CreatedAt)
	u, _ := ptypes.Timestamp(role.UpdatedAt)

	return Role{
		ID:        role.Id,
		UUID:      role.Uuid,
		Title:     role.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
