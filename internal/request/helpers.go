package request

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func getStringFrom(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func getSliceStringFrom(values []string) []string {
	if len(values) > 0 {
		return strings.Split(values[0], ",")
	}

	return []string{}
}

func getSliceArrayFrom(values []string) []string {
	var req []string
	if len(values) > 0 {
		req = append(req, values...)
	}

	return req
}

func getOptionalStringFrom(values []string) *string {
	if len(values) > 0 {
		return &values[0]
	}
	return nil
}

func getOptionalInt64From(values []string) *int64 {
	if len(values) == 0 || values[0] == "" {
		return nil
	}
	val, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil
	}
	return &val
}

func getUintFrom(values []string) uint {
	if len(values) > 0 {
		v, _ := strconv.ParseUint(values[0], 10, 32)
		return uint(v)
	}
	return 0
}

func getFloat64From(values []string) float64 {
	if len(values) > 0 {
		v, _ := strconv.ParseFloat(values[0], 64)
		return v
	}
	return 0
}

func getIntFrom(values []string) int {
	if len(values) > 0 {
		v, _ := strconv.ParseInt(values[0], 10, 64)
		return int(v)
	}
	return 0
}

func getBoolFrom(values []string) bool {
	if len(values) > 0 {
		v, _ := strconv.ParseBool(values[0])
		return v
	}
	return false
}

func getIntFromQuery(value string, defaultValue int) int {
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return defaultValue
}

func getStringOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func stringPtr(s string) *string {
	return &s
}

func getFormValue(values map[string][]string, key string) string {
	if v, ok := values[key]; ok && len(v) > 0 {
		return v[0]
	}
	return ""
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func parseDate(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)
	return t
}

func parseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func parseFloatPtr(s string) *float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return nil
	}
	return &f
}

func parseStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}

// func GetValidationError(err error, originalStruct interface{}) map[string][]string {
// 	originalValue := reflect.ValueOf(originalStruct)
// 	validatorErr := map[string][]string{}

// 	if validationErrors, ok := err.(validator.ValidationErrors); ok {
// 		for _, e := range validationErrors {
// 			fieldName := e.StructField()
// 			label := getLabelFromTag(originalValue, fieldName)
// 			message := e.Translate(validatorPkg.Trans)
// 			if label != "" {
// 				message = strings.Replace(message, e.Field(), label, 1)
// 			}

// 			validatorErr[e.Field()] = append(validatorErr[e.Field()], message)
// 		}
// 	}

// 	return validatorErr
// }

func getLabelFromTag(originalValue reflect.Value, fieldName string) string {
	// Ensure we have the original struct and not a pointer to it
	if originalValue.Kind() == reflect.Ptr {
		originalValue = originalValue.Elem()
	}

	// Get the type of the original struct
	structType := originalValue.Type()

	// Find the field by name
	if field, found := structType.FieldByName(fieldName); found {
		// Return the value of the "label" tag if it exists
		label := field.Tag.Get("label")
		if label != "" {
			return label
		}

		tag := field.Tag.Get("json")
		if tag != "" {
			return tag
		}
	}

	return "-"
}
