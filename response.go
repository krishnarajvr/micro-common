package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response defines the api response
type Response struct {
	Status    int         `json:"status" example:"200"`
	Data      interface{} `json:"data,omitempty" example:"{data:{products}}"`
	Error     interface{} `json:"error,omitempty" example:"{}"`
	RequestId string      `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ResponseWithPage struct {
	Status     int                    `json:"status" example:"200"`
	Data       map[string]interface{} `json:"data,omitempty" example:"{data:{products}}"`
	Error      interface{}            `json:"error,omitempty" example:"{}"`
	Pagination interface{}            `json:"_pagination,omitempty" example:"{}"`
	RequestId  string                 `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

func SuccessResponse(c *gin.Context, key string, body interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusOK,
		Data:      map[string]interface{}{key: body},
		Error:     nil,
		RequestId: c.Request.Header.Get("X-B3-Traceid"),
	})
}

func SuccessPageResponse(c *gin.Context, key string, body interface{}, page interface{}) {
	dataWithPage := ResponseWithPage{
		Status:     http.StatusOK,
		Data:       map[string]interface{}{key: body},
		Error:      nil,
		Pagination: page,
		RequestId:  c.Request.Header.Get("X-B3-Traceid"),
	}

	c.JSON(http.StatusOK, dataWithPage)
}

func ErrorResponseWitCode(c *gin.Context, errorCode int, errorData *ErrorData) {
	c.JSON(errorCode, Response{
		Status:    errorCode,
		Data:      nil,
		Error:     errorData,
		RequestId: c.Request.Header.Get("X-B3-Traceid"),
	})
}

func ErrorResponse(c *gin.Context, errorData *ErrorData) {
	c.JSON(http.StatusBadRequest, Response{
		Status:    http.StatusBadRequest,
		Data:      nil,
		Error:     errorData,
		RequestId: c.Request.Header.Get("X-B3-Traceid"),
	})
}

func BadRequest(c *gin.Context, errorData interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Status:    http.StatusBadRequest,
		Data:      nil,
		Error:     errorData,
		RequestId: c.Request.Header.Get("X-B3-Traceid"),
	})
}

func BadRequestWithMessage(c *gin.Context, errorMessage string) {
	if len(errorMessage) == 0 {
		errorMessage = "Bad Request"
	}

	errorData := &ErrorData{
		Code:    BAD_REQUEST,
		Message: errorMessage,
	}

	ErrorResponseWitCode(c, http.StatusBadRequest, errorData)
}

func AccessDenied(c *gin.Context, errorMessage string) {
	if len(errorMessage) == 0 {
		errorMessage = "Access Denied"
	}

	errorData := &ErrorData{
		Code:    ACCESS_DENIED,
		Message: errorMessage,
	}

	ErrorResponseWitCode(c, http.StatusForbidden, errorData)
}

func ResourceNotFound(c *gin.Context, errorMessage string) {
	if len(errorMessage) == 0 {
		errorMessage = "Requested resource not found"
	}

	errorData := &ErrorData{
		Code:    NOT_FOUND,
		Message: errorMessage,
	}

	ErrorResponseWitCode(c, http.StatusNotFound, errorData)
}

func InternalServerError(c *gin.Context, errorMessage string) {
	if len(errorMessage) == 0 {
		errorMessage = "Internal server error"
	}

	errorData := &ErrorData{
		Code:    INTERNAL_SERVER_ERROR,
		Message: errorMessage,
	}

	ErrorResponseWitCode(c, http.StatusInternalServerError, errorData)
}

func PartialSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusPartialContent, Response{
		Status:    http.StatusPartialContent,
		Data:      data,
		RequestId: c.Request.Header.Get("X-B3-Traceid"),
	})
}
