package db

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"

	//dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/mirzakhany/pm/pkg/config"
	"github.com/mirzakhany/pm/pkg/log"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *pg.DB
}

// TransactionFunc represents a function that will start a transaction and run the given function.
type TransactionFunc func(ctx context.Context, f func(ctx context.Context) error) error

type contextKey int

const (
	txKey contextKey = iota
)

var (
	host     config.String
	port     config.Int
	db       config.String
	password config.String
	user     config.String
)

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *pg.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(ctx context.Context) *pg.DB {
	if tx, ok := ctx.Value(txKey).(*pg.DB); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

func Init(ctx context.Context) (*DB, error) {

	host = config.RegisterString("db.host", "localhost")
	port = config.RegisterInt("db.port", 5432)
	db = config.RegisterString("db.db", "postgres")
	password = config.RegisterString("db.password", "postgres")
	user = config.RegisterString("db.user", "postgres")

	if err := config.Load(); err != nil {
		log.Panic("load database settings failed")
		return nil, err
	}

	database := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", host.String(), port.Int()),
		User:     user.String(),
		Password: password.String(),
		Database: db.String(),
	})

	go func() {
		<-ctx.Done()
		if err := database.Close(); err != nil {
			log.Error("error in close database", log.Err(err))
		}
	}()

	return &DB{database}, nil
}
