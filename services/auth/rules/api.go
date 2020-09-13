package rules

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

// RuleRequest create or update rule request
type RuleRequest struct {
	Subject  string `form:"subject" json:"subject" xml:"subject" binding:"required"`
	Domain   string `form:"domain" json:"domain" xml:"domain" binding:"required"`
	Resource string `form:"resource" json:"resource" xml:"resource" binding:"required"`
	Action   string `form:"action" json:"action" xml:"action" binding:"required"`
	Object   string `form:"object" json:"object" xml:"object" binding:"required"`
}

// RuleResponse create or update rule request
type RuleResponse struct {
	UUID     string `json:"uuid" xml:"uuid"`
	Subject  string `form:"subject" json:"subject" xml:"subject"`
	Domain   string `form:"domain" json:"domain" xml:"domain"`
	Resource string `form:"resource" json:"resource" xml:"resource"`
	Action   string `form:"action" json:"action" xml:"action"`
	Object   string `form:"object" json:"object" xml:"object"`
}

func NewService(router *gin.Engine, repo *Repository, logger *zap.Logger) {

	service := &Service{
		Repo:        repo,
		Logger:      logger,
		itemPerPage: 10,
	}

	router.POST("/rules", service.CreateHandler)
	router.PUT("/rules/:uuid", service.UpdateHandler)
	router.GET("/rules", service.RetrieveHandler)
	router.GET("/rules/:uuid", service.GetHandler)
	router.DELETE("/rules/:uuid", service.DeleteHandler)
}

func (s *Service) CreateHandler(c *gin.Context) {
	var json RuleRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := s.Repo.Create(json)
	if err != nil {
		s.Logger.Error("create rule failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) UpdateHandler(c *gin.Context) {
	var json RuleRequest
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
		s.Logger.Error("update rule failed with error %s", zap.Error(err))
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

		s.Logger.Error("retrieve rule failed with error %s", zap.Error(err))
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
		s.Logger.Error("retrieve rules failed with error %s", zap.Error(err))
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
		s.Logger.Error("delete rule failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
