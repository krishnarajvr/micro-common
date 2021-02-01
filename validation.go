package common

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetMessages - Get the message based on the language
func GetMessages(langKey string) map[string]string {
	return map[string]string{
		"Required":     "Can not be empty",
		"Min":          "Minimum is %d",
		"Max":          "Maximum is %d",
		"Range":        "Range is %d to %d",
		"MinSize":      "Minimum size is %d",
		"MaxSize":      "Maximum size is %d",
		"Length":       "Required length is %d",
		"Alpha":        "Must be valid alpha characters",
		"Numeric":      "Must be valid numeric characters",
		"AlphaNumeric": "Must be valid alpha or numeric characters",
		"Match":        "Must match %s",
		"NoMatch":      "Must not match %s",
		"AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
		"Email":        "Must be a valid email address",
		"IP":           "Must be a valid ip address",
		"Base64":       "Must be valid base64 characters",
		"Mobile":       "Must be valid mobile number",
		"Tel":          "Must be valid telephone number",
		"Phone":        "Must be valid telephone or mobile phone number",
		"ZipCode":      "Must be valid zipcode",
	}
}

// ValidateForm - Validate the data from request based on the form definition
func ValidateForm(c *gin.Context, form interface{}) (status bool, errorData *ErrorData) {
	err := c.Bind(form)
	if err != nil {
		return true, nil
	}

	log := c.MustGet("log").(*logrus.Logger)

	validation.MessageTmpls = GetMessages("en")
	valid := validation.Validation{}
	check, err := valid.Valid(form)

	if err != nil {
		errorData = &ErrorData{
			Code:    INTERNAL_SERVER_ERROR,
			Message: "Internal Server Error",
		}
		return false, errorData
	}

	if !check {
		log.Info("Validation error")
		log.Info(check)
		log.Info(valid.Errors)
		errorDetails := make([]ErrorDetail, 0)

		for _, err := range valid.Errors {
			log.Info(err.Key, err.Message)
			errorDetails = append(errorDetails,
				ErrorDetail{
					Code:    err.Key,
					Target:  err.Field,
					Message: err.Message,
				},
			)
		}

		errorData = &ErrorData{
			Code:    BAD_REQUEST,
			Message: "BAD Request",
			Details: errorDetails,
		}

		return false, errorData
	}

	return true, nil
}
