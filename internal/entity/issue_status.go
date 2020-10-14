package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/issues"
)

type IssueStatus struct {
	tableName struct{} `pg:"issues_status,alias:ss"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ss IssueStatus) ToProto(secure bool) *issues.IssueStatus {
	c, _ := ptypes.TimestampProto(ss.CreatedAt)
	u, _ := ptypes.TimestampProto(ss.UpdatedAt)

	issueStatus := &issues.IssueStatus{
		Id:        ss.ID,
		Uuid:      ss.UUID,
		Title:     ss.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return issueStatus
}

func IssueStatusToProtoList(iml []IssueStatus, secure bool) []*issues.IssueStatus {
	var _issueStatus []*issues.IssueStatus
	for _, i := range iml {
		_issueStatus = append(_issueStatus, i.ToProto(secure))
	}
	return _issueStatus
}

func IssueStatusFromProto(issueStatus *issues.IssueStatus) IssueStatus {
	c, _ := ptypes.Timestamp(issueStatus.CreatedAt)
	u, _ := ptypes.Timestamp(issueStatus.UpdatedAt)

	return IssueStatus{
		ID:        issueStatus.Id,
		UUID:      issueStatus.Uuid,
		Title:     issueStatus.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
