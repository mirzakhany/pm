package db

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mirzakhany/pm/pkg/config"
	"testing"
)

var testDB *DB

var (
	testDSN config.String
)

// NewForTest returns the database connection for testing purpose.
func NewForTest(t *testing.T, models []interface{}) *DB {
	if testDB != nil {
		return testDB
	}

	testDSN = config.RegisterString("db.testDSN", "postgres://postgres:postgres@localhost:5432/task_manager_test?sslmode=disable")
	opt, err := pg.ParseURL(testDSN.String())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dbc := pg.Connect(opt)

	err = createSchema(dbc, models)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return &DB{dbc}
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *DB, tables ...string) {
	for _, table := range tables {
		_, err := db.DB().Exec(fmt.Sprintf("TRUNCATE %s", table))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

// createSchema creates database schema for test
func createSchema(db *pg.DB, models []interface{}) error {
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
