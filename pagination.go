package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination - Pagination request data from API
type Pagination struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
	Search string `json:"search"`
}

// PageResult - Pagination output to API
type PageResult struct {
	PageNumber int   `json:"pageNumber"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int64 `json:"totalPages"`
	IsFirst    int   `json:"isFirst"`
	IsLast     int   `json:"isLast"`
}

//Paginator - Populate pagination object from request query
func Paginator(c *gin.Context) Pagination {
	var page Pagination

	page.Sort = c.DefaultQuery("sort", "ID")
	page.Order = c.DefaultQuery("order", "DESC")
	page.Offset = c.DefaultQuery("offset", "0")
	page.Limit = c.DefaultQuery("limit", "25")
	page.Search = c.DefaultQuery("Search", "")

	return page
}

// PageInfo - Gives the paginatio info to API
func PageInfo(page Pagination, totalCount int64) PageResult {
	var pageResult PageResult

	pageResult.TotalCount = totalCount
	pageResult.PageSize, _ = strconv.Atoi(page.Limit)

	totalPages := totalCount / int64(pageResult.PageSize)
	pageResult.TotalPages = totalPages

	pageResult.IsFirst = 1
	pageResult.IsLast = 0

	offset, _ := strconv.Atoi(page.Offset)

	if offset > 1 {
		pageResult.IsFirst = 0
	}

	if offset >= int(pageResult.TotalPages) {
		pageResult.IsLast = 1
	}

	return pageResult
}

//Paginate - Do the db pagination query
func Paginate(page Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		limitNumber, _ := strconv.Atoi(page.Limit)
		offsetNumber, _ := strconv.Atoi(page.Offset)

		offset := (offsetNumber - 1) * limitNumber
		return db.Offset(offset).Limit(limitNumber).Order(page.Sort + " " + page.Order)
	}
}
