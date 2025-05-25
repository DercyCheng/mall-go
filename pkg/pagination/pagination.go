// Package pagination provides utilities for handling pagination in API responses
package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params represents pagination parameters from the request
type Params struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Offset   int `json:"offset,omitempty"`
	Limit    int `json:"limit,omitempty"`
}

// Response represents a paginated response
type Response struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	HasMore    bool        `json:"hasMore"`
}

// DefaultParams returns the default pagination parameters
func DefaultParams() *Params {
	return &Params{
		Page:     1,
		PageSize: 10,
	}
}

// NewParams creates a new Params object from the given values
func NewParams(page, pageSize int) *Params {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		// Limit maximum page size to prevent resource exhaustion
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	return &Params{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// FromGinContext extracts pagination parameters from the Gin context
func FromGinContext(c *gin.Context) *Params {
	// Get pagination parameters from query
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// Parse parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		// Limit maximum page size to prevent resource exhaustion
		pageSize = 100
	}

	return NewParams(page, pageSize)
}

// NewResponse creates a new paginated response
func NewResponse(items interface{}, total int64, params *Params) *Response {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	hasMore := params.Page < totalPages

	return &Response{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}
