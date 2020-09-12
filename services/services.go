package services

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/services/auth"
	"projectmanager/services/projects"
)

func Setup(db *gorm.DB, router *gin.Engine, logger *zap.Logger) error {

	err := projects.Setup(db, router, logger)
	if err != nil {
		return err
	}

	err = auth.Setup(db, router, logger)
	return err
}
