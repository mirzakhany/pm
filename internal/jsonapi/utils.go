package jsonapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PaginationQuery(c *gin.Context, perPage int) (int, int, error) {
	offset := c.DefaultQuery("offset", "0")
	limit := c.DefaultQuery("limit", strconv.Itoa(perPage))

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid offset number %s", offset)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid limit number %s", limit)
	}

	return offsetInt, limitInt, nil
}
