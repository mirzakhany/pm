package server

import (
	"projectmanager/services"
	"time"

	"moul.io/zapgorm2"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Start() error {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbLogger := zapgorm2.New(logger)
	dbLogger.SetAsDefault()

	dsn := "host=localhost user=postgres password=postgres dbname=task_manager port=5432 sslmode=disable TimeZone=Europe/Stockholm"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{Logger: dbLogger})
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	err = services.Setup(db, r, logger)
	if err != nil {
		return err
	}

	// Listen and Server in 0.0.0.0:8080
	return r.Run(":8080")
}
