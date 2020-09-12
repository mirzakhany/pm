package jsonapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Single(status int, c *gin.Context, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

func List(status int, c *gin.Context, data interface{}, offset, limit int, total int64) {

	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}

	nextOffset := offset + limit

	links := gin.H{"self": c.FullPath()}

	if offset > 0 {
		links["prev"] = fmt.Sprintf("%s?offset=%d&limit=%d", c.FullPath(), prevOffset, limit)
	}

	if int64(nextOffset) < total {
		links["next"] = fmt.Sprintf("%s?offset=%d&limit=%d", c.FullPath(), offset+limit, limit)
	}

	c.JSON(status, gin.H{
		"pagination": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
		},
		"data":  data,
		"links": links,
	})
}
