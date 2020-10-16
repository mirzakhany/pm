package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/workspaces"
)

type Workspace struct {
	tableName struct{} `pg:"workspaces,alias:w"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Title     string
	Domain    string `pg:",unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (rm Workspace) ToProto() *workspaces.Workspace {
	c, _ := ptypes.TimestampProto(rm.CreatedAt)
	u, _ := ptypes.TimestampProto(rm.UpdatedAt)

	workspace := &workspaces.Workspace{
		Id:        rm.ID,
		Uuid:      rm.UUID,
		Title:     rm.Title,
		Domain:    rm.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return workspace
}

func WorkspaceToProtoList(rml []Workspace) []*workspaces.Workspace {
	var r []*workspaces.Workspace
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func WorkspaceFromProto(workspace *workspaces.Workspace) Workspace {
	c, _ := ptypes.Timestamp(workspace.CreatedAt)
	u, _ := ptypes.Timestamp(workspace.UpdatedAt)

	return Workspace{
		ID:        workspace.Id,
		UUID:      workspace.Uuid,
		Title:     workspace.Title,
		Domain:    workspace.Domain,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
