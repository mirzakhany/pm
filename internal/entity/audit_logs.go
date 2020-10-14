package entity

import "time"

type AuditLog struct {
	tableName struct{} `pg:"audit_logs,alias:al"` //nolint
	ID        uint64
	Action    string
	Object    string
	OldData   string
	NewData   string
	ByID      uint64 `pg:"unique:by_id"`
	By        *User  `pg:"rel:has-one"`
	CreatedAt time.Time
}
