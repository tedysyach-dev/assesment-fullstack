package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Common error constructors
func NewBadRequestError(message string, details interface{}) *AppError {
	return &AppError{
		Code:    fiber.StatusBadRequest,
		Message: message,
		Details: details,
	}
}

func NewValidationError(err error) *AppError {
	validationErrors := ParseValidationErrors(err)
	return &AppError{
		Code:     fiber.StatusBadRequest,
		Message:  "Validation failed",
		Details:  validationErrors,
		Internal: err,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    fiber.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    fiber.StatusUnauthorized,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:     fiber.StatusInternalServerError,
		Message:  "Internal server error",
		Internal: err,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    fiber.StatusConflict,
		Message: message,
	}
}

func ParseValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if validatorErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validatorErrs {
			fieldName := getJSONFieldName(e)

			validationErrors = append(validationErrors, ValidationError{
				Field:   fieldName,
				Message: getValidationMessage(e, fieldName),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
			})
		}
	}

	return validationErrors
}

func getJSONFieldName(err validator.FieldError) string {
	field := err.Field()

	if field == "" {
		return toSnakeCase(err.StructField())
	}

	return field
}

func getValidationMessage(err validator.FieldError, field string) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	var prevLower bool

	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			if prevLower {
				result.WriteRune('_')
			} else if i+1 < len(s) {
				nextRune := rune(s[i+1])
				if nextRune >= 'a' && nextRune <= 'z' {
					result.WriteRune('_')
				}
			}
		}

		prevLower = (r >= 'a' && r <= 'z')
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}
