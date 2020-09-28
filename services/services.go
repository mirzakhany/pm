package services

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mirzakhany/pm/pkg/db"
	rolesSrv "github.com/mirzakhany/pm/services/roles"
	tasksSrv "github.com/mirzakhany/pm/services/tasks"
	usersSrv "github.com/mirzakhany/pm/services/users"
)

func Setup(db *db.DB) error {

	err := createSchema(db.DB())
	if err != nil {
		return err
	}

	rolesSrv.New(rolesSrv.NewService(rolesSrv.NewRepository(db)))
	userService := usersSrv.NewService(usersSrv.NewRepository(db))
	usersSrv.New(userService)
	tasksSrv.New(tasksSrv.NewService(tasksSrv.NewRepository(db), userService))
	return nil
}

// createSchema creates database schema
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*usersSrv.UserModel)(nil),
		(*rolesSrv.RoleModel)(nil),
		(*tasksSrv.TaskModel)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
