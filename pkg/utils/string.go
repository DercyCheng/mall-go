// Package utils provides common utility functions
package utils

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// String conversion utilities

// ToString converts any value to a string
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case []byte:
		return string(v)
	default:
		// Try JSON marshaling for complex types
		if jsonBytes, err := json.Marshal(value); err == nil {
			return string(jsonBytes)
		}
		// Fall back to reflect for other types
		return reflect.ValueOf(value).String()
	}
}

// ToInt converts a value to an int with a default value if conversion fails
func ToInt(value interface{}, defaultValue int) int {
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	case bool:
		if v {
			return 1
		}
		return 0
	}

	return defaultValue
}

// String manipulation utilities

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Truncate truncates a string to a specified length with an optional suffix
func Truncate(s string, maxLength int, suffix string) string {
	if len(s) <= maxLength {
		return s
	}

	return s[:maxLength] + suffix
}

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}

// SecureRandomString generates a secure random string of the specified length
func SecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// CamelToSnake converts a camelCase string to snake_case
func CamelToSnake(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SnakeToCamel converts a snake_case string to camelCase
func SnakeToCamel(s string) string {
	var result strings.Builder

	nextUpper := false
	for _, r := range s {
		if r == '_' {
			nextUpper = true
		} else {
			if nextUpper {
				result.WriteRune(unicode.ToUpper(r))
				nextUpper = false
			} else {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}

// Map utilities

// MapGet safely gets a value from a map with a default value if the key doesn't exist
func MapGet(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// MergeMaps merges multiple maps into a single map
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
