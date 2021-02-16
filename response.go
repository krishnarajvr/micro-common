package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response defines the api response
type Response struct {
	Status int         `json:"code" example:"200"`
	Data   interface{} `json:"data" example:"{data:{products}}"`
	Error  interface{} `json:"error" example:"{}"`
}

type ResponseWithPage struct {
	Status     int                    `json:"status" example:"200"`
	Data       map[string]interface{} `json:"data" example:"{data:{products}}"`
	Error      interface{}            `json:"error" example:"{}"`
	Pagination interface{}            `json:"_pagination" example:"{}"`
}

func SuccessResponse(c *gin.Context, key string, body interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   map[string]interface{}{key: body},
		"error":  "",
	})
}

func SuccessPageResponse(c *gin.Context, key string, body interface{}, page interface{}) {
	dataWithPage := ResponseWithPage{
		Status:     http.StatusOK,
		Data:       map[string]interface{}{key: body},
		Error:      "",
		Pagination: page,
	}

	c.JSON(http.StatusOK, dataWithPage)
}

func BadRequest(c *gin.Context, errorData interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"data":   "",
		"error":  errorData,
	})
}

func ErrorResponse(c *gin.Context, errorData *ErrorData) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"data":   "",
		"error":  errorData,
	})
}
