package internal

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	cyclesSrv "github.com/mirzakhany/pm/internal/cycles"
	"github.com/mirzakhany/pm/internal/entity"
	issuesSrv "github.com/mirzakhany/pm/internal/issues"
	rolesSrv "github.com/mirzakhany/pm/internal/roles"
	usersSrv "github.com/mirzakhany/pm/internal/users"
	"github.com/mirzakhany/pm/pkg/db"
)

func Setup(db *db.DB) error {

	err := createSchema(db.DB())
	if err != nil {
		return err
	}

	rolesSrv.New(rolesSrv.NewService(rolesSrv.NewRepository(db)))
	userService := usersSrv.NewService(usersSrv.NewRepository(db))
	usersSrv.New(userService)
	cycleService := cyclesSrv.NewService(cyclesSrv.NewRepository(db), userService)
	cyclesSrv.New(cycleService)
	issuesSrv.New(issuesSrv.NewService(issuesSrv.NewRepository(db), userService, cycleService))
	return nil
}

// createSchema creates database schema
func createSchema(db *pg.DB) error {
	models := []interface{}{
		&entity.User{},
		&entity.Cycle{},
		&entity.Role{},
		&entity.IssueStatus{},
		&entity.Issue{},
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
