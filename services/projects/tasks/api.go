package tasks

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"projectmanager/internal/jsonapi"
)

type Service struct {
	Repo        *Repository
	Logger      *zap.Logger
	itemPerPage int
}

// TaskRequest create or update task request
type TaskRequest struct {
	Title    string `form:"title" json:"title" xml:"title" binding:"required"`
	SprintID int    `form:"sprintId" json:"sprintId" xml:"sprint" binding:"required"`
	Estimate string `form:"estimate" json:"estimate" xml:"estimate" binding:"required"`
	Status   string `form:"status" json:"status" xml:"status" binding:"required"`
	Assignee string `form:"assignee" json:"assignee" xml:"assignee" binding:"required"`
}

// TaskRequest create or update task request
type TaskResponse struct {
	UUID     string `json:"uuid" xml:"uuid"`
	Title    string `json:"title" xml:"title"`
	Sprint   string `json:"sprintId" xml:"sprintID"`
	Estimate string `json:"estimate" xml:"estimate"`
	Status   string `json:"status" xml:"status"`
	Assignee string `json:"assignee" xml:"assignee"`
}

func NewService(router *gin.Engine, repo *Repository, logger *zap.Logger) {

	service := &Service{
		Repo:        repo,
		Logger:      logger,
		itemPerPage: 20,
	}

	router.POST("/tasks", service.CreateHandler)
	router.PUT("/tasks/:uuid", service.UpdateHandler)
	router.GET("/tasks", service.RetrieveHandler)
	router.GET("/tasks/:uuid", service.GetHandler)
	router.DELETE("/tasks/:uuid", service.DeleteHandler)
}

func (s *Service) CreateHandler(c *gin.Context) {
	var json TaskRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := s.Repo.Create(json)
	if err != nil {
		s.Logger.Error("create task failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonapi.Single(http.StatusOK, c, item)
}

func (s *Service) UpdateHandler(c *gin.Context) {
	var json TaskRequest
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
		s.Logger.Error("update item failed with error %s", zap.Error(err))
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
		s.Logger.Error("retrieve task failed with error %s", zap.Error(err))
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

	items, totalCount, err := s.Repo.GetTasks(offset, limit)
	if err != nil {
		s.Logger.Error("retrieve tasks failed with error %s", zap.Error(err))
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
		s.Logger.Error("delete task failed with error %s", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
