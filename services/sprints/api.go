package sprints

import (
	"errors"
	"net/http"
	"projectmanager/internal/jsonapi"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Repo        *Repository
	Logger      *zap.Logger
	itemPerPage int
}

// SprintRequest create or update sprint request
type SprintRequest struct {
	Title   string    `form:"title" json:"title" xml:"title" binding:"required"`
	Status  string    `form:"status" json:"status" xml:"status" binding:"required"`
	StartAt time.Time `form:"startAt" json:"startAt" xml:"startAt" binding:"required"`
	EndAt   time.Time `form:"endAt" json:"endAt" xml:"endAt" binding:"required"`
}

// SprintRequest create or update task request
type SprintResponse struct {
	UUID    string    `json:"uuid" xml:"uuid"`
	Title   string    `json:"title" xml:"title"`
	Status  string    `json:"status" xml:"status" binding:"required"`
	StartAt time.Time `json:"startAt" xml:"startAt" binding:"required"`
	EndAt   time.Time `json:"endAt" xml:"endAt" binding:"required"`
}

func NewService(router *gin.Engine, repo *Repository, logger *zap.Logger) {

	service := &Service{
		Repo:        repo,
		Logger:      logger,
		itemPerPage: 10,
	}

	router.POST("/sprints", service.CreateHandler)
	router.PUT("/sprints/:uuid", service.UpdateHandler)
	router.GET("/sprints", service.RetrieveHandler)
	router.GET("/sprints/:uuid", service.GetHandler)
	router.DELETE("/sprints/:uuid", service.DeleteHandler)
}

func (s *Service) CreateHandler(c *gin.Context) {
	var json SprintRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := s.Repo.Create(json)
	if err != nil {
		s.Logger.Error("create sprint failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) UpdateHandler(c *gin.Context) {
	var json SprintRequest
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
		s.Logger.Error("update sprint failed with error %s", zap.Error(err))
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

		s.Logger.Error("retrieve sprint failed with error %s", zap.Error(err))
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
		s.Logger.Error("retrieve sprints failed with error %s", zap.Error(err))
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
	err := s.Repo.DeleteSprint(itemUUID)
	if err != nil {
		s.Logger.Error("delete sprint failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
