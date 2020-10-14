package entity

import "time"

type Workspace struct {
	ID        uint64 `pg:"workspaces,alias:w"` //nolint
	UUID      string `pg:"default:gen_random_uuid()"`
	Name      string
	Domain    string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
