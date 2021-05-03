package common

import (
	"regexp"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xeipuuv/gojsonschema"
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
		"Binary":       "Must be 0 or 1",
		"PhoneNo":      "Must be valid telephone or mobile phone number",
		"PositiveInt":  "Must be valid positive integer",
	}
}

// ValidateForm - Validate the data from request based on the form definition
func ValidateForm(c *gin.Context, form interface{}) (status bool, errorData *ErrorData) {
	err := c.Bind(form)
	if err != nil {
		return true, nil
	}

	log := c.MustGet("log").(*MicroLog)

	validation.MessageTmpls = GetMessages("en")
	valid := validation.Validation{}
	validation.AddCustomFunc("Binary", Binary)
	validation.AddCustomFunc("PhoneNo", PhoneNo)
	validation.AddCustomFunc("PositiveInt", PositiveInt)
	check, err := valid.Valid(form)

	if err != nil {
		errorData = &ErrorData{
			Code:    INTERNAL_SERVER_ERROR,
			Message: "Internal Server Error",
		}
		return false, errorData
	}

	if !check {
		log.Message("Validation error")
		log.Message(valid.Errors)
		errorDetails := make([]ErrorDetail, 0)

		for _, err := range valid.Errors {
			log.Message(err.Key + ":" + err.Message)
			Name := err.Key
			Field := ""
			Label := ""
			parts := strings.Split(err.Key, ".")
			if len(parts) == 3 {
				Field = parts[0]
				Name = parts[1]
				Label = parts[2]
				if len(Label) == 0 {
					Label = Field
				}
			}
			errorDetails = append(errorDetails,
				ErrorDetail{
					Code:    Name,
					Target:  Label,
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

// ValidateID - Validate ID fields
func ValidateID(c *gin.Context, key string) (int, bool) {
	const MIN_ID_VALUE int = 1
	id := com.StrTo(c.Param(key)).MustInt()
	valid := validation.Validation{}
	valid.Min(id, MIN_ID_VALUE, key)

	if valid.HasErrors() {
		return 0, false
	}

	return id, true
}

// Custom JSON schema validation
func ValidateCustomSchema(metaSchemaDef *map[string]interface{}, form *map[string]interface{}, validationConfig map[string]string) (*map[string]interface{}, *ErrorData) {
	schemaRequiredFields := (*metaSchemaDef)["required"].([]interface{})
	var errorData *ErrorData

	if validationConfig["trimRequired"] != "false" {
		form = TrimMapFields(form, schemaRequiredFields)
	}

	result, err := gojsonschema.Validate(gojsonschema.NewGoLoader(&metaSchemaDef), gojsonschema.NewGoLoader(&form))

	if err != nil {
		errorData = &ErrorData{
			Code:    INTERNAL_SERVER_ERROR,
			Message: "Internal Server Error",
		}

		return nil, errorData
	}

	if !result.Valid() {
		errorDetails := make([]ErrorDetail, 0)

		for _, err := range result.Errors() {
			var errMap = err.Details()
			var target interface{}

			if errMap["field"] == "(root)" {
				target = errMap["property"]

			} else {
				target = errMap["field"]
			}

			errorDetails = append(errorDetails,
				ErrorDetail{
					// Further enchancement required for Code
					Code:    "",
					Target:  target.(string),
					Message: err.Description(),
				},
			)
		}

		errorData = &ErrorData{
			Code:    BAD_REQUEST,
			Message: "BAD Request",
			Details: errorDetails,
		}

		return nil, errorData
	}

	return form, nil
}

// Binary is custom form field validation for binary type 0,1
func Binary(v *validation.Validation, obj interface{}, key string) {
	binary, ok := obj.(int)

	if !ok {
		return
	}

	messages := GetMessages("en")
	if binary != 0 && binary != 1 {
		v.AddError(key, messages["Binary"])
		return
	}
}

// Phone is custom form field validation for phone no
func PhoneNo(v *validation.Validation, obj interface{}, key string) {
	phoneNo := obj.(string)
	//Further enhancement required based on country specific pattern
	var mobilePattern = regexp.MustCompile(`^[0-9+()-]*$`)

	matchphone := mobilePattern.MatchString(phoneNo)
	if matchphone {
		return
	} else {
		messages := GetMessages("en")
		v.AddError(key, messages["PhoneNo"])
		return

	}
}

// Trim spaces for required fields.
func TrimMapFields(form *map[string]interface{}, fields []interface{}) *map[string]interface{} {
	for _, requiredField := range fields {
		if (*form)[requiredField.(string)] != nil {
			fieldTrimmed := strings.TrimSpace((*form)[requiredField.(string)].(string))
			(*form)[requiredField.(string)] = fieldTrimmed
		}
	}
	return form
}

// ValidateCode - Validate Code fields
func ValidateCode(c *gin.Context, key string) (string, bool) {
	code := com.StrTo(c.Param(key)).String()
	codeValid := validation.Validation{}
	codeValid.AlphaDash(code, key)

	if codeValid.HasErrors() {
		return code, false
	}

	return code, true
}

// PositiveInt is custom form field validation for int data type for validate positive integer 1,2,..n
func PositiveInt(v *validation.Validation, obj interface{}, key string) {
	messages := GetMessages("en")
	positiveInt, ok := obj.(int)

	if !ok {
		v.AddError(key, messages["PositiveInt"])
		return
	}

	if positiveInt <= 0 {
		v.AddError(key, messages["PositiveInt"])
		return
	}
}
