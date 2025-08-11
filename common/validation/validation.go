package validation

import (
"fmt"
"regexp"
"strings"
"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
Field   string `json:"field"`
Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
var messages []string
for _, err := range ve {
messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
}
return strings.Join(messages, ", ")
}

// Validator provides validation functions
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
return &Validator{}
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
if !emailRegex.MatchString(email) {
return fmt.Errorf("invalid email format")
}
return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
if len(password) < 6 {
return fmt.Errorf("password must be at least 6 characters long")
}

hasUpper := false
hasLower := false
hasNumber := false

for _, char := range password {
if unicode.IsUpper(char) {
hasUpper = true
} else if unicode.IsLower(char) {
hasLower = true
} else if unicode.IsNumber(char) {
hasNumber = true
}
}

if !hasUpper || !hasLower || !hasNumber {
return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one number")
}

return nil
}

// ValidateRequired validates that a field is not empty
func (v *Validator) ValidateRequired(value string, fieldName string) error {
if strings.TrimSpace(value) == "" {
return fmt.Errorf("%s is required", fieldName)
}
return nil
}

// ValidateMinLength validates minimum length
func (v *Validator) ValidateMinLength(value string, minLength int, fieldName string) error {
if len(value) < minLength {
return fmt.Errorf("%s must be at least %d characters long", fieldName, minLength)
}
return nil
}

// ValidateMaxLength validates maximum length
func (v *Validator) ValidateMaxLength(value string, maxLength int, fieldName string) error {
if len(value) > maxLength {
return fmt.Errorf("%s must be no more than %d characters long", fieldName, maxLength)
}
return nil
}

// ValidateUsername validates username format
func (v *Validator) ValidateUsername(username string) error {
if err := v.ValidateRequired(username, "username"); err != nil {
return err
}

if err := v.ValidateMinLength(username, 3, "username"); err != nil {
return err
}

if err := v.ValidateMaxLength(username, 50, "username"); err != nil {
return err
}

// Check if username contains only alphanumeric characters
for _, char := range username {
if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
return fmt.Errorf("username can only contain letters and numbers")
}
}

return nil
}

// ValidateSlug validates slug format
func (v *Validator) ValidateSlug(slug string) error {
if slug == "" {
return nil // Slug can be empty, it will be generated
}

slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
if !slugRegex.MatchString(slug) {
return fmt.Errorf("slug can only contain lowercase letters, numbers, and hyphens")
}

if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
return fmt.Errorf("slug cannot start or end with a hyphen")
}

if strings.Contains(slug, "--") {
return fmt.Errorf("slug cannot contain consecutive hyphens")
}

return nil
}

// ValidateStatus validates post/page status
func (v *Validator) ValidateStatus(status string) error {
validStatuses := []string{"draft", "published"}
for _, validStatus := range validStatuses {
if status == validStatus {
return nil
}
}
return fmt.Errorf("status must be one of: %s", strings.Join(validStatuses, ", "))
}

// ValidateRole validates user role
func (v *Validator) ValidateRole(role string) error {
if role == "" {
return nil // Role can be empty, default will be used
}

validRoles := []string{"admin", "editor", "user"}
for _, validRole := range validRoles {
if role == validRole {
return nil
}
}
return fmt.Errorf("role must be one of: %s", strings.Join(validRoles, ", "))
}

// ValidateID validates that an ID is positive
func (v *Validator) ValidateID(id int, fieldName string) error {
if id <= 0 {
return fmt.Errorf("%s must be a positive integer", fieldName)
}
return nil
}
