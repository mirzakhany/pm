package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/issues"
)

type Issue struct {
	tableName   struct{} `pg:"issues,alias:i"` //nolint
	ID          uint64   `pg:",pk"`
	UUID        string   `pg:"default:gen_random_uuid()"`
	Title       string
	Description string
	StatusID    uint64       `pg:"unique:status_id"`
	Status      *IssueStatus `pg:"rel:has-one, fk:status"`
	CycleID     uint64       `pg:"unique:cycle_id"`
	Cycle       *Cycle       `pg:"rel:has-one, fk:cycle"`
	Estimate    uint64
	AssigneeID  uint64 `pg:"unique:assignee_id"`
	Assignee    *User  `pg:"rel:has-one, fk:assignee"`
	CreatorID   uint64 `pg:"unique:creator_id"`
	Creator     *User  `pg:"rel:has-one, fk:creator"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (im Issue) ToProto(secure bool) *issues.Issue {
	c, _ := ptypes.TimestampProto(im.CreatedAt)
	u, _ := ptypes.TimestampProto(im.UpdatedAt)

	status := im.Status.ToProto(false)
	cycle := &issues.Issue{
		Id:          im.ID,
		Uuid:        im.UUID,
		Title:       im.Title,
		Description: im.Description,
		Status:      status,
		Cycle:       im.Cycle.ToProto(secure),
		Estimate:    im.Estimate,
		Creator:     im.Creator.ToProto(secure),
		Assignee:    im.Assignee.ToProto(secure),
		CreatedAt:   c,
		UpdatedAt:   u,
	}
	return cycle
}

func IssueToProtoList(iml []Issue, secure bool) []*issues.Issue {
	var _issues []*issues.Issue
	for _, i := range iml {
		_issues = append(_issues, i.ToProto(secure))
	}
	return _issues
}

func IssueFromProto(issue *issues.Issue) Issue {
	c, _ := ptypes.Timestamp(issue.CreatedAt)
	u, _ := ptypes.Timestamp(issue.UpdatedAt)

	assignee := UserFromProto(issue.Creator)
	creator := UserFromProto(issue.Creator)
	cycle := CycleFromProto(issue.Cycle)
	issueStatus := IssueStatusFromProto(issue.Status)
	return Issue{
		ID:          issue.Id,
		UUID:        issue.Uuid,
		Title:       issue.Title,
		Description: issue.Description,
		CycleID:     cycle.ID,
		Cycle:       &cycle,
		Creator:     &creator,
		Status:      &issueStatus,
		StatusID:    issueStatus.ID,
		CreatorID:   creator.ID,
		Assignee:    &assignee,
		AssigneeID:  assignee.ID,
		CreatedAt:   c,
		UpdatedAt:   u,
	}
}
