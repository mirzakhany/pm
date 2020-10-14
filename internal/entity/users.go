package entity

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mirzakhany/pm/protobuf/users"
)

type User struct {
	tableName struct{} `pg:"users,alias:u"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Username  string   `pg:",unique"`
	Password  string
	Email     string `pg:",unique"`
	Enable    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (um User) ToProto(secure bool) *users.User {
	c, _ := ptypes.TimestampProto(um.CreatedAt)
	u, _ := ptypes.TimestampProto(um.UpdatedAt)

	user := &users.User{
		Id:        um.ID,
		Uuid:      um.UUID,
		Username:  um.Username,
		Password:  um.Password,
		Email:     um.Email,
		Enable:    um.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	if secure {
		user.Id = 0
		user.Password = ""
	}
	return user
}

func UserToProtoList(uml []User, secure bool) []*users.User {
	var u []*users.User
	for _, i := range uml {
		u = append(u, i.ToProto(secure))
	}
	return u
}

func UserFromProto(user *users.User) User {
	c, _ := ptypes.Timestamp(user.CreatedAt)
	u, _ := ptypes.Timestamp(user.UpdatedAt)
	return User{
		ID:        user.Id,
		UUID:      user.Uuid,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Enable:    user.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
