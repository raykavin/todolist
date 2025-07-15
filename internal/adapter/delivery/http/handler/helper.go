package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"todolist/internal/adapter/delivery/http"

	"github.com/go-playground/validator/v10"
)

func getIDParam(ctx http.RequestContext) (int64, error) {
	id, err := strconv.ParseInt(ctx.GetParam("id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter")
	}

	return id, nil
}

// getAuthenticatedUserID extracts and validates the user ID from the context
func getAuthenticatedUserID(ctx http.RequestContext) (int64, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("invalid user ID type")
	}

	return id, nil
}

// parseError parses gin binding errors into a map
func parseError(err error) map[string]any {
	errors := make(map[string]any)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				if e.Type().Kind() == reflect.String {
					errors[field] = fmt.Sprintf("Must be at least %s characters", e.Param())
				} else {
					errors[field] = fmt.Sprintf("Must be at least %s", e.Param())
				}
			case "max":
				if e.Type().Kind() == reflect.String {
					errors[field] = fmt.Sprintf("Must not exceed %s characters", e.Param())
				} else {
					errors[field] = fmt.Sprintf("Must not exceed %s", e.Param())
				}
			case "oneof":
				errors[field] = fmt.Sprintf("Must be one of: %s", e.Param())
			default:
				errors[field] = fmt.Sprintf("Failed on %s validation", tag)
			}
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}
