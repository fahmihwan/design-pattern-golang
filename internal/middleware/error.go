package middleware

import (
	util "best-pattern/internal/util/errorUtil"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func HandleValidationErrors(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "application/json")

	// Handle validator package errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		validationMap := make(map[string][]string)
		for _, verr := range validationErrors {
			field := convertFieldName(verr.Field())
			validationMap[field] = append(validationMap[field], verr.Tag())
		}
		util.SendErrorResponse(w, http.StatusBadRequest, validationMap)
		return true
	}

	// Handle invalid validation errors
	if _, ok := err.(*validator.InvalidValidationError); ok {
		util.SendErrorResponse(w, http.StatusInternalServerError, map[string]string{"general": err.Error()})
		return true
	}

	// Handle regular errors
	util.SendErrorResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	return true
}

// Convert CamelCase to snake_case
func convertFieldName(fieldName string) string {
	var result []rune
	for i, r := range fieldName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
