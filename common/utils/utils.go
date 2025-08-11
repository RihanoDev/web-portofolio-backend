package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Pagination struct for paginated queries
type Pagination struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// NewPagination creates a new pagination struct with validated values
func NewPagination(page, limit int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// StringToSlug converts a string to URL-friendly slug
func StringToSlug(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)

	var result strings.Builder
	lastWasHyphen := false

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
			lastWasHyphen = false
		} else if !lastWasHyphen && result.Len() > 0 {
			result.WriteRune('-')
			lastWasHyphen = true
		}
	}

	// Remove trailing hyphen
	slug := result.String()
	if strings.HasSuffix(slug, "-") {
		slug = slug[:len(slug)-1]
	}

	return slug
}

// Truncate truncates string to specified length
func Truncate(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[:length] + "..."
}

// IsValidEmail validates email format (basic validation)
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ParseID converts string to uint ID
func ParseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %v", err)
	}
	return uint(id), nil
}

// ParseIntID parses string ID to int
func ParseIntID(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}

// ValidatePageAndLimit validates and normalizes pagination parameters
func ValidatePageAndLimit(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 { // Maximum limit
		limit = 100
	}
	return page, limit
}
