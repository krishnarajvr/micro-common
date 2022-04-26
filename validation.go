package common

import (
	"regexp"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xeipuuv/gojsonschema"
)

func init() {
	// Moved to init function because of : (fatal error: concurrent map writes) during performance testing
	validation.AddCustomFunc("Binary", Binary)
	validation.AddCustomFunc("PhoneNo", PhoneNo)
	validation.AddCustomFunc("PositiveInt", PositiveInt)
	validation.AddCustomFunc("BinaryPointer", BinaryPointer)
}

type JsonSchemaLocale struct{}

// False returns a format-string for "false" schema validation errors
func (l JsonSchemaLocale) False() string {
	return "False always fails validation"
}

// Required returns a format-string for "required" schema validation errors
func (l JsonSchemaLocale) Required() string {
	return `{{.property}} is required`
}

// InvalidType returns a format-string for "invalid type" schema validation errors
func (l JsonSchemaLocale) InvalidType() string {
	return `Invalid type. Expected: {{.expected}}, given: {{.given}}`
}

// NumberAnyOf returns a format-string for "anyOf" schema validation errors
func (l JsonSchemaLocale) NumberAnyOf() string {
	return `Must validate at least one schema (anyOf)`
}

// NumberOneOf returns a format-string for "oneOf" schema validation errors
func (l JsonSchemaLocale) NumberOneOf() string {
	return `Must validate one and only one schema (oneOf)`
}

// NumberAllOf returns a format-string for "allOf" schema validation errors
func (l JsonSchemaLocale) NumberAllOf() string {
	return `Must validate all the schemas (allOf)`
}

// NumberNot returns a format-string to format a NumberNotError
func (l JsonSchemaLocale) NumberNot() string {
	return `Must not validate the schema (not)`
}

// MissingDependency returns a format-string for "missing dependency" schema validation errors
func (l JsonSchemaLocale) MissingDependency() string {
	return `Has a dependency on {{.dependency}}`
}

// Internal returns a format-string for internal errors
func (l JsonSchemaLocale) Internal() string {
	return `Internal Error {{.error}}`
}

// Const returns a format-string to format a ConstError
func (l JsonSchemaLocale) Const() string {
	return `{{.field}} does not match: {{.allowed}}`
}

// Enum returns a format-string to format an EnumError
func (l JsonSchemaLocale) Enum() string {
	return `{{.field}} must be one of the following: {{.allowed}}`
}

// ArrayNoAdditionalItems returns a format-string to format an ArrayNoAdditionalItemsError
func (l JsonSchemaLocale) ArrayNoAdditionalItems() string {
	return `No additional items allowed on array`
}

// ArrayNotEnoughItems returns a format-string to format an error for arrays having not enough items to match positional list of schema
func (l JsonSchemaLocale) ArrayNotEnoughItems() string {
	return `Not enough items on array to match positional list of schema`
}

// ArrayMinItems returns a format-string to format an ArrayMinItemsError
func (l JsonSchemaLocale) ArrayMinItems() string {
	return `Array must have at least {{.min}} items`
}

// ArrayMaxItems returns a format-string to format an ArrayMaxItemsError
func (l JsonSchemaLocale) ArrayMaxItems() string {
	return `Array must have at most {{.max}} items`
}

// Unique returns a format-string  to format an ItemsMustBeUniqueError
func (l JsonSchemaLocale) Unique() string {
	return `{{.type}} items[{{.i}},{{.j}}] must be unique`
}

// ArrayContains returns a format-string to format an ArrayContainsError
func (l JsonSchemaLocale) ArrayContains() string {
	return `At least one of the items must match`
}

// ArrayMinProperties returns a format-string to format an ArrayMinPropertiesError
func (l JsonSchemaLocale) ArrayMinProperties() string {
	return `Must have at least {{.min}} properties`
}

// ArrayMaxProperties returns a format-string to format an ArrayMaxPropertiesError
func (l JsonSchemaLocale) ArrayMaxProperties() string {
	return `Must have at most {{.max}} properties`
}

// AdditionalPropertyNotAllowed returns a format-string to format an AdditionalPropertyNotAllowedError
func (l JsonSchemaLocale) AdditionalPropertyNotAllowed() string {
	return `Additional property {{.property}} is not allowed`
}

// InvalidPropertyPattern returns a format-string to format an InvalidPropertyPatternError
func (l JsonSchemaLocale) InvalidPropertyPattern() string {
	return `Property "{{.property}}" does not match pattern {{.pattern}}`
}

// InvalidPropertyName returns a format-string to format an InvalidPropertyNameError
func (l JsonSchemaLocale) InvalidPropertyName() string {
	return `Property name of "{{.property}}" does not match`
}

// StringGTE returns a format-string to format an StringLengthGTEError
func (l JsonSchemaLocale) StringGTE() string {
	return `String length must be greater than or equal to {{.min}}`
}

// StringLTE returns a format-string to format an StringLengthLTEError
func (l JsonSchemaLocale) StringLTE() string {
	return `String length must be less than or equal to {{.max}}`
}

// DoesNotMatchPattern returns a format-string to format an DoesNotMatchPatternError
func (l JsonSchemaLocale) DoesNotMatchPattern() string {
	return `Does not match pattern '{{.pattern}}'`
}

// DoesNotMatchFormat returns a format-string to format an DoesNotMatchFormatError
func (l JsonSchemaLocale) DoesNotMatchFormat() string {
	return `Does not match format '{{.format}}'`
}

// MultipleOf returns a format-string to format an MultipleOfError
func (l JsonSchemaLocale) MultipleOf() string {
	return `Must be a multiple of {{.multiple}}`
}

// NumberGTE returns the format string to format a NumberGTEError
func (l JsonSchemaLocale) NumberGTE() string {
	return `Must be greater than or equal to {{.min}}`
}

// NumberGT returns the format string to format a NumberGTError
func (l JsonSchemaLocale) NumberGT() string {
	return `Must be greater than {{.min}}`
}

// NumberLTE returns the format string to format a NumberLTEError
func (l JsonSchemaLocale) NumberLTE() string {
	return `Must be less than or equal to {{.max}}`
}

// NumberLT returns the format string to format a NumberLTError
func (l JsonSchemaLocale) NumberLT() string {
	return `Must be less than {{.max}}`
}

// Schema validators

// RegexPattern returns a format-string to format a regex-pattern error
func (l JsonSchemaLocale) RegexPattern() string {
	return `Invalid regex pattern '{{.pattern}}'`
}

// GreaterThanZero returns a format-string to format an error where a number must be greater than zero
func (l JsonSchemaLocale) GreaterThanZero() string {
	return `{{.number}} must be strictly greater than 0`
}

// MustBeOfA returns a format-string to format an error where a value is of the wrong type
func (l JsonSchemaLocale) MustBeOfA() string {
	return `{{.x}} must be of a {{.y}}`
}

// MustBeOfAn returns a format-string to format an error where a value is of the wrong type
func (l JsonSchemaLocale) MustBeOfAn() string {
	return `{{.x}} must be of an {{.y}}`
}

// CannotBeUsedWithout returns a format-string to format a "cannot be used without" error
func (l JsonSchemaLocale) CannotBeUsedWithout() string {
	return `{{.x}} cannot be used without {{.y}}`
}

// CannotBeGT returns a format-string to format an error where a value are greater than allowed
func (l JsonSchemaLocale) CannotBeGT() string {
	return `{{.x}} cannot be greater than {{.y}}`
}

// MustBeOfType returns a format-string to format an error where a value does not match the required type
func (l JsonSchemaLocale) MustBeOfType() string {
	return `{{.key}} must be of type {{.type}}`
}

// MustBeValidRegex returns a format-string to format an error where a regex is invalid
func (l JsonSchemaLocale) MustBeValidRegex() string {
	return `{{.key}} must be a valid regex`
}

// MustBeValidFormat returns a format-string to format an error where a value does not match the expected format
func (l JsonSchemaLocale) MustBeValidFormat() string {
	return `{{.key}} must be a valid format {{.given}}`
}

// MustBeGTEZero returns a format-string to format an error where a value must be greater or equal than 0
func (l JsonSchemaLocale) MustBeGTEZero() string {
	return `{{.key}} must be greater than or equal to 0`
}

// KeyCannotBeGreaterThan returns a format-string to format an error where a value is greater than the maximum  allowed
func (l JsonSchemaLocale) KeyCannotBeGreaterThan() string {
	return `{{.key}} cannot be greater than {{.y}}`
}

// KeyItemsMustBeOfType returns a format-string to format an error where a key is of the wrong type
func (l JsonSchemaLocale) KeyItemsMustBeOfType() string {
	return `{{.key}} items must be {{.type}}`
}

// KeyItemsMustBeUnique returns a format-string to format an error where keys are not unique
func (l JsonSchemaLocale) KeyItemsMustBeUnique() string {
	return `{{.key}} items must be unique`
}

// ReferenceMustBeCanonical returns a format-string to format a "reference must be canonical" error
func (l JsonSchemaLocale) ReferenceMustBeCanonical() string {
	return `Reference {{.reference}} must be canonical`
}

// NotAValidType returns a format-string to format an invalid type error
func (l JsonSchemaLocale) NotAValidType() string {
	return `has a primitive type that is NOT VALID -- given: {{.given}} Expected valid values are:{{.expected}}`
}

// Duplicated returns a format-string to format an error where types are duplicated
func (l JsonSchemaLocale) Duplicated() string {
	return `{{.type}} type is duplicated`
}

// HttpBadStatus returns a format-string for errors when loading a schema using HTTP
func (l JsonSchemaLocale) HttpBadStatus() string {
	return `Could not read schema from HTTP, response status is {{.status}}`
}

// ErrorFormat returns a format string for errors
// Replacement options: field, description, context, value
func (l JsonSchemaLocale) ErrorFormat() string {
	return `{{.field}}: {{.description}}`
}

// ParseError returns a format-string for JSON parsing errors
func (l JsonSchemaLocale) ParseError() string {
	return `Expected: {{.expected}}, given: Invalid JSON`
}

// ConditionThen returns a format-string for ConditionThenError errors
// If/Else
func (l JsonSchemaLocale) ConditionThen() string {
	return `Must validate "then" as "if" was valid`
}

// ConditionElse returns a format-string for ConditionElseError errors
func (l JsonSchemaLocale) ConditionElse() string {
	return `Must validate "else" as "if" was not valid`
}

// GetMessages - Get the message based on the language
func GetMessages(langKey string) map[string]string {
	return map[string]string{
		"Required":      "Can not be empty",
		"Min":           "Minimum is %d",
		"Max":           "Maximum is %d",
		"Range":         "Range is %d to %d",
		"MinSize":       "Minimum size is %d",
		"MaxSize":       "Maximum size is %d",
		"Length":        "Required length is %d",
		"Alpha":         "Must be valid alpha characters",
		"Numeric":       "Must be valid numeric characters",
		"AlphaNumeric":  "Must be valid alpha or numeric characters",
		"Match":         "Must match %s",
		"NoMatch":       "Must not match %s",
		"AlphaDash":     "Must be valid alpha or numeric or dash(-_) characters",
		"Email":         "Must be a valid email address",
		"IP":            "Must be a valid ip address",
		"Base64":        "Must be valid base64 characters",
		"Mobile":        "Must be valid mobile number",
		"Tel":           "Must be valid telephone number",
		"Phone":         "Must be valid telephone or mobile phone number",
		"ZipCode":       "Must be valid zipcode",
		"Binary":        "Must be 0 or 1",
		"PhoneNo":       "Must be valid telephone or mobile phone number",
		"PositiveInt":   "Must be valid positive integer",
		"BinaryPointer": "Must be 0 or 1",
	}
}

// SimpleValidateForm - Validate the data from request based on the form definition (Without AddCustomFunc support)
func SimpleValidateForm(c *gin.Context, form interface{}) (status bool, errorData *ErrorData) {
	err := c.Bind(form)
	if err != nil {
		return true, nil
	}

	log := c.MustGet("log").(*MicroLog)

	validation.MessageTmpls = GetMessages("en")
	valid := validation.Validation{}
	check, err := valid.RecursiveValid(form)

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

// ValidateForm - Validate the data from request based on the form definition
func ValidateForm(c *gin.Context, form interface{}) (status bool, errorData *ErrorData) {
	err := c.Bind(form)
	if err != nil {
		return true, nil
	}

	log := c.MustGet("log").(*MicroLog)

	validation.MessageTmpls = GetMessages("en")
	valid := validation.Validation{}
	check, err := valid.RecursiveValid(form)

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

	gojsonschema.Locale = JsonSchemaLocale{}
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
					Code:    err.Type(),
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

// BinaryPointer is custom form field validation for binary int type 0,1
func BinaryPointer(v *validation.Validation, obj interface{}, key string) {
	messages := GetMessages("en")
	binary, ok := obj.(*int)

	if !ok {
		v.AddError(key, "Data type should be int pointer")
		return
	}

	if binary != nil && *binary != 0 && *binary != 1 {
		v.AddError(key, messages["BinaryPointer"])
		return
	}
}
