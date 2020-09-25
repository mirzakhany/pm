package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"proj/pkg/config"
	"proj/pkg/log"
)

var (
	engine   config.String
	host     config.String
	port     config.Int
	db       config.String
	password config.String
	user     config.String
	timeZone config.String
)

func postgresEngine() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host.String(), user.String(), password.String(), db.String(), port.Int(), timeZone.String(),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: log.DbLogger()})
	if err != nil {
		log.Fatal("unable to connect to postgres database", log.Err(err))
		return nil, err
	}
	return db, err
}

func Init() (*gorm.DB, error) {

	engine = config.RegisterString("db.engine", "postgres")
	host = config.RegisterString("db.host", "localhost")
	port = config.RegisterInt("db.port", 5432)
	db = config.RegisterString("db.db", "postgres")
	password = config.RegisterString("db.password", "postgres")
	user = config.RegisterString("db.user", "postgres")
	timeZone = config.RegisterString("db.timeZone", "Europe/Stockholm")

	if err := config.Load(); err != nil {
		log.Panic("load database settings failed")
		return nil, err
	}

	if engine.String() == "postgres" {
		return postgresEngine()
	}

	db, err := gorm.Open(sqlite.Open(host.String()), &gorm.Config{Logger: log.DbLogger()})
	if err != nil {
		log.Fatal("unable to connect to sqlite database", log.Err(err))
		return nil, err
	}
	return db, nil
}
