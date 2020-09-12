package projects

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/services/projects/sprints"
	"projectmanager/services/projects/tasks"
)

func Setup(db *gorm.DB, router *gin.Engine, logger *zap.Logger) error {

	err := db.AutoMigrate(&sprints.Sprint{}, &tasks.Task{})
	if err != nil {
		logger.Fatal("migrate projects tables failed", zap.Error(err))
		return err
	}

	tasks.NewService(router, tasks.NewRepository(db), logger)
	sprints.NewService(router, sprints.NewRepository(db), logger)

	return nil
}
