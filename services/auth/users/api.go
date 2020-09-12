package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"projectmanager/internal/jsonapi"
)

type Service struct {
	Repo        *Repository
	Logger      *zap.Logger
	itemPerPage int
}

// UserRequest create or update user request
type UserRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	Email    string `form:"email" json:"email" xml:"email" binding:"required"`
	Enable   bool   `form:"enable" json:"enable" xml:"enable" binding:"required"`
}

// UserResponse create or update user request
type UserResponse struct {
	UUID     string `json:"uuid" xml:"uuid"`
	Username string `form:"username" json:"username" xml:"username"`
	Password string `form:"password" json:"password" xml:"password"`
	Email    string `form:"email" json:"email" xml:"email"`
	Enable   bool   `form:"enable" json:"enable" xml:"enable"`
}

func NewService(router *gin.Engine, repo *Repository, logger *zap.Logger) {

	service := &Service{
		Repo:        repo,
		Logger:      logger,
		itemPerPage: 10,
	}

	router.POST("/users", service.CreateHandler)
	router.PUT("/users/:uuid", service.UpdateHandler)
	router.GET("/users", service.RetrieveHandler)
	router.GET("/users/:uuid", service.GetHandler)
	router.DELETE("/users/:uuid", service.DeleteHandler)
}

func (s *Service) CreateHandler(c *gin.Context) {
	var json UserRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := s.Repo.Create(json)
	if err != nil {
		s.Logger.Error("create user failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) UpdateHandler(c *gin.Context) {
	var json UserRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemUUID := c.Param("uuid")
	if itemUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid not provided"})
		return
	}

	item, err := s.Repo.Update(itemUUID, json)
	if err != nil {
		s.Logger.Error("update user failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) GetHandler(c *gin.Context) {

	itemUUID := c.Param("uuid")
	if itemUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid not provided"})
		return
	}

	item, err := s.Repo.GetByUUID(itemUUID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		s.Logger.Error("retrieve user failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) RetrieveHandler(c *gin.Context) {

	offset, limit, err := jsonapi.PaginationQuery(c, s.itemPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, totalCount, err := s.Repo.Retrieve(offset, limit)
	if err != nil {
		s.Logger.Error("retrieve users failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.List(http.StatusOK, c, items, offset, limit, totalCount)
}

func (s *Service) DeleteHandler(c *gin.Context) {
	itemUUID := c.Param("uuid")
	if itemUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid not provided"})
		return
	}
	err := s.Repo.Delete(itemUUID)
	if err != nil {
		s.Logger.Error("delete user failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
