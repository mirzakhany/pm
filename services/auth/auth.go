package auth

import (
	"projectmanager/services/auth/domains"
	"projectmanager/services/auth/roles"
	"projectmanager/services/auth/rules"
	"projectmanager/services/auth/users"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, router *gin.Engine, logger *zap.Logger) error {

	err := db.AutoMigrate(&users.User{}, roles.Role{}, rules.Rule{}, domains.Domain{})
	if err != nil {
		logger.Fatal("migrate auth tables failed", zap.Error(err))
		return err
	}
	users.NewService(router, users.NewRepository(db), logger)
	roles.NewService(router, roles.NewRepository(db), logger)
	rules.NewService(router, rules.NewRepository(db), logger)
	domains.NewService(router, domains.NewRepository(db), logger)
	return nil
}
