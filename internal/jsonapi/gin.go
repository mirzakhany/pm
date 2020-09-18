package jsonapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type PaginationData struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total"`
}

type LinksData struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
	Self string `json:"self"`
}

type PaginatedRes struct {
	Pagination PaginationData `json:"pagination"`
	Data       interface{}    `json:"data"`
	Links      LinksData      `json:"links"`
}

type SingleRes struct {
	Data interface{} `json:"data"`
}

func Single(status int, c *gin.Context, data interface{}) {
	c.JSON(status, SingleRes{Data: data})
}

func List(status int, c *gin.Context, data interface{}, offset, limit int, total int64) {

	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}

	nextOffset := offset + limit

	res := PaginatedRes{
		Data: data,
		Pagination: PaginationData{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
		Links: LinksData{
			Self: c.FullPath(),
		},
	}

	if offset > 0 {
		res.Links.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", c.FullPath(), prevOffset, limit)
	}

	if int64(nextOffset) < total {
		res.Links.Next = fmt.Sprintf("%s?offset=%d&limit=%d", c.FullPath(), offset+limit, limit)
	}

	c.JSON(status, res)
}
