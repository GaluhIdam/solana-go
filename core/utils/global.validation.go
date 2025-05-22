// utils/validation.go
package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GlobalValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func getValidationMessage(tag, field, param string) string {
	field = capitalize(field)
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case "ne":
		return fmt.Sprintf("%s cannot be equal to %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "boolean":
		return fmt.Sprintf("%s must be a boolean", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime", field)
	case "contains":
		return fmt.Sprintf("%s must contain %s", field, param)
	case "excludes":
		return fmt.Sprintf("%s cannot contain %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, param)
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", field)
	case "ipv4":
		return fmt.Sprintf("%s must be a valid IPv4 address", field)
	case "ipv6":
		return fmt.Sprintf("%s must be a valid IPv6 address", field)
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}

func MapValidationErrors(err error) []GlobalValidation {
	var validations []GlobalValidation
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			tag := e.Tag()
			param := e.Param()
			message := getValidationMessage(tag, field, param)
			message = strings.ToLower(message)

			validations = append(validations, GlobalValidation{
				Field:   field,
				Message: message,
			})
		}
	} else {
		validations = append(validations, GlobalValidation{
			Field:   "",
			Message: err.Error(),
		})
	}
	return validations
}

func MustBindAndValidate(c *fiber.Ctx, out interface{}, validateFunc func() error) {
	if err := c.BodyParser(out); err != nil {
		panic(fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	if err := validateFunc(); err != nil {
		panic(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Validation Error",
			"errors":  MapValidationErrors(err),
		})
	}
}
