package roles

import (
	"errors"
	"net/http"
	"projectmanager/internal/jsonapi"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Repo        *Repository
	Logger      *zap.Logger
	itemPerPage int
}

// RoleRequest create or update role request
type RoleRequest struct {
	Title string `form:"title" json:"title" xml:"title" binding:"required"`
}

// RoleResponse create or update role request
type RoleResponse struct {
	UUID  string `json:"uuid" xml:"uuid"`
	Title string `form:"title" json:"title" xml:"title"`
}

func NewService(router *gin.Engine, repo *Repository, logger *zap.Logger) {

	service := &Service{
		Repo:        repo,
		Logger:      logger,
		itemPerPage: 10,
	}

	router.POST("/roles", service.CreateHandler)
	router.PUT("/roles/:uuid", service.UpdateHandler)
	router.GET("/roles", service.RetrieveHandler)
	router.GET("/roles/:uuid", service.GetHandler)
	router.DELETE("/roles/:uuid", service.DeleteHandler)
}

func (s *Service) CreateHandler(c *gin.Context) {
	var json RoleRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := s.Repo.Create(json)
	if err != nil {
		s.Logger.Error("create role failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) UpdateHandler(c *gin.Context) {
	var json RoleRequest
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
		s.Logger.Error("update role failed with error %s", zap.Error(err))
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

		s.Logger.Error("retrieve role failed with error %s", zap.Error(err))
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
		s.Logger.Error("retrieve roles failed with error %s", zap.Error(err))
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
		s.Logger.Error("delete role failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
