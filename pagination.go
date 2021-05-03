package common

import (
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

// Pagination - Pagination request data from API
type Pagination struct {
	PageNumber string `json:"pageNumber"`
	//pageOffset will take precedence over pageNumber if provided
	PageOffset string `json:"pageOffset"`
	PageSize   string `json:"pageSize"`
	PageOrder  string `json:"pageOrder"`
	Search     string `json:"q"`
}

//Page - Valid page object created from Pagination request
type Page struct {
	Number int
	Offset int
	Size   int
	Order  string
	Search string
}

// PageResult - Pagination output to API
type PageResult struct {
	PageNumber int   `json:"pageNumber"`
	PageOffset int   `json:"pageOffset"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int64 `json:"totalPages"`
	IsFirst    int   `json:"isFirst"`
	IsLast     int   `json:"isLast"`
}

//Paginator - Populate pagination object from request query
func Paginator(c *gin.Context) Pagination {
	var page Pagination

	page.PageNumber = c.DefaultQuery("pageNumber", "1")
	page.PageSize = c.DefaultQuery("pageSize", "25")
	page.PageOrder = c.DefaultQuery("pageOrder", "")
	page.PageOffset = c.DefaultQuery("pageOffset", "")

	return page
}

// PageInfo - Gives the paginatio info to API
func PageInfo(pagination Pagination, totalCount int64) PageResult {
	var pageResult PageResult

	page := ValidPage(pagination, map[string]interface{}{})

	pageResult.TotalCount = totalCount
	pageResult.PageSize = page.Size
	pageResult.PageNumber = page.Number
	pageResult.PageOffset = page.Offset

	totalPages := float64(totalCount) / float64(pageResult.PageSize)
	pageResult.TotalPages = int64(math.Ceil(totalPages))

	pageResult.IsFirst = 1
	pageResult.IsLast = 0

	if pageResult.PageNumber > 1 {
		pageResult.IsFirst = 0
	}

	if (pageResult.PageOffset + pageResult.PageSize) > int(totalCount) {
		pageResult.IsLast = 1
	}

	return pageResult
}

//ValidPage return the page object after validation
func ValidPage(pagination Pagination, allowedFields map[string]interface{}) Page {
	var page Page

	page.Offset, _ = strconv.Atoi(pagination.PageOffset)
	page.Size, _ = strconv.Atoi(pagination.PageSize)
	page.Number, _ = strconv.Atoi(pagination.PageNumber)
	page.Order = ""

	//Todo: Get configuration from external
	if page.Size <= 0 {
		page.Size = 25
	}

	if page.Number <= 0 {
		page.Number = 1
	}

	if page.Offset > 0 {
		page.Number = ((page.Offset - 1) / page.Size) + 1
	}

	//Todo: Get configuration from external
	if page.Size > 500 {
		page.Size = 500
	}

	if page.Offset <= 0 {
		page.Offset = ((page.Number - 1) * page.Size) + 1
	}

	page.Order = parseOrder(pagination.PageOrder, allowedFields)

	return page
}

func parseOrder(order string, allowedFields map[string]interface{}) string {
	var orderBy string
	var orderFieldWithSort []string
	var orderFieldSnake string

	if len(allowedFields) == 0 {
		allowedFields = map[string]interface{}{
			"id":         "true",
			"name":       "true",
			"code":       "true",
			"created_at": "true",
			"updated_at": "true",
		}
	}

	orderString := ""

	if len(order) == 0 {
		return orderString
	}

	orders := strings.Split(order, ",")

	for _, orderField := range orders {
		orderField = strings.TrimSpace(orderField)
		orderFieldWithSort = strings.Split(orderField, " ")

		if len(orderFieldWithSort) >= 1 {
			if len(orderFieldWithSort) == 1 {
				orderBy = "DESC"
			} else {
				orderBy = orderFieldWithSort[1]
				orderBy = strings.TrimSpace(orderBy)
				orderBy = strings.ToUpper(orderBy)
			}

			orderField = strings.TrimSpace(orderFieldWithSort[0])
			orderFieldSnake = strcase.ToSnake(orderField)

			if orderBy != "ASC" && orderBy != "DESC" {
				orderBy = "DESC"
			}

			if _, ok := allowedFields[orderFieldSnake]; ok {
				if orderString == "" {
					orderString = orderFieldSnake + " " + orderBy
				} else {
					orderString = orderString + ", " + orderFieldSnake + " " + orderBy
				}
			}
		}
	}

	return orderString
}

//Paginate - Do the db pagination query
func Paginate(pagination Pagination) func(db *gorm.DB) *gorm.DB {
	page := ValidPage(pagination, map[string]interface{}{})

	return func(db *gorm.DB) *gorm.DB {
		db = db.Offset(page.Offset - 1).Limit(page.Size)
		if len(page.Order) > 0 {
			db = db.Order(page.Order)
		}

		return db
	}
}
