package templates

const EntityTmpl = `
package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/{{ .Pkg.NamePlural | lower }}"
)

type {{ .Pkg.Name }} struct {
	tableName struct{} ` + "`pg:\"{{ .Pkg.NamePlural | lower }},alias:{{ .Pkg.EntityAlias }}\"`" + ` //nolint
	ID        uint64   ` + "`pg:\",pk\"`" + `
	UUID      string   ` + "`pg:\"default:gen_random_uuid()\"`" + `
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (rm {{ .Pkg.Name }}) ToProto() *{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }} {
	c, _ := ptypes.TimestampProto(rm.CreatedAt)
	u, _ := ptypes.TimestampProto(rm.UpdatedAt)

	{{ .Pkg.Name | lower }} := &{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}{
		Id:        rm.ID,
		Uuid:      rm.UUID,
		Title:     rm.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return {{ .Pkg.Name | lower }}
}

func {{ .Pkg.Name }}ToProtoList(rml []{{ .Pkg.Name }}) []*{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }} {
	var r []*{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func {{ .Pkg.Name }}FromProto({{ .Pkg.Name | lower }} *{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}) {{ .Pkg.Name }} {
	c, _ := ptypes.Timestamp({{ .Pkg.Name | lower }}.CreatedAt)
	u, _ := ptypes.Timestamp({{ .Pkg.Name | lower }}.UpdatedAt)

	return {{ .Pkg.Name }}{
		ID:        {{ .Pkg.Name | lower }}.Id,
		UUID:      {{ .Pkg.Name | lower }}.Uuid,
		Title:     {{ .Pkg.Name | lower }}.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
}

`
