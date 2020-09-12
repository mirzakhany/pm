package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/services/auth/roles"
	"projectmanager/services/auth/rules"
	"projectmanager/services/auth/users"
)

func Setup(db *gorm.DB, router *gin.Engine, logger *zap.Logger) error {

	err := db.AutoMigrate(&users.User{}, roles.Role{}, rules.Rule{})
	if err != nil {
		logger.Fatal("migrate auth tables failed", zap.Error(err))
		return err
	}
	users.NewService(router, users.NewRepository(db), logger)
	roles.NewService(router, roles.NewRepository(db), logger)
	rules.NewService(router, rules.NewRepository(db), logger)
	return nil
}
