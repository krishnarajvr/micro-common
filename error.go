package common

import "strings"

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

const (
	BAD_REQUEST              = "BAD_REQUEST"
	INTERNAL_SERVER_ERROR    = "INTERNAL_SERVER_ERROR"
	INVALID_ARGUMENT         = "INVALID_ARGUMENT"
	OUT_OF_RANGE             = "OUT_OF_RANGE"
	UNAUTHENTICATED          = "UNAUTHENTICATED"
	ACCESS_DENIED            = "ACCESS_DENIED"
	NOT_FOUND                = "NOT_FOUND"
	ABORTED                  = "ABORTED"
	ALREADY_EXISTS           = "ALREADY_EXISTS"
	RESOURCE_EXHAUSTED       = "RESOURCE_EXHAUSTED"
	CANCELLED                = "CANCELLED"
	DATA_LOSS                = "DATA_LOSS"
	UNKNOWN                  = "UNKNOWN"
	NOT_IMPLEMENTED          = "NOT_IMPLEMENTED"
	UNAVAILABLE              = "UNAVAILABLE"
	DEADLINE_EXCEEDED        = "DEADLINE_EXCEEDED"
	REFERENCE_INTEGRITY_FAIL = "REFERENCE_INTEGRITY_FAIL"
)

type ErrorData struct {
	Code    string        `json:"code" example:"BAD_REQUEST"`
	Message string        `json:"message" example:"Bad Request"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code" example:"Required"`
	Target  string `json:"target" example:"Name"`
	Message string `json:"message" example:"Name field is required"`
}

func CheckDbError(err error) string {
	if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
		return ALREADY_EXISTS
	}
	if strings.HasPrefix(err.Error(), "Error 1452: Cannot add or update a child") {
		return REFERENCE_INTEGRITY_FAIL
	}

	if strings.HasPrefix(err.Error(), "record not found") {
		return NOT_FOUND
	}

	return INTERNAL_SERVER_ERROR
}
