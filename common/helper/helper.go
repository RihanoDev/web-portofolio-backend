package helper

import (
"crypto/rand"
"fmt"
"html"
"regexp"
"strings"
"time"
)

// Helper provides utility functions
type Helper struct{}

// NewHelper creates a new helper instance
func NewHelper() *Helper {
return &Helper{}
}

// FormatDate formats a time.Time to a readable string
func (h *Helper) FormatDate(t time.Time, layout string) string {
if layout == "" {
layout = "2006-01-02 15:04:05"
}
return t.Format(layout)
}

// FormatDateHuman formats a time.Time to human-readable format
func (h *Helper) FormatDateHuman(t time.Time) string {
return t.Format("January 2, 2006")
}

// SanitizeHTML removes HTML tags from a string
func (h *Helper) SanitizeHTML(input string) string {
// Remove HTML tags
re := regexp.MustCompile(`<[^>]*>`)
text := re.ReplaceAllString(input, "")

// Unescape HTML entities
text = html.UnescapeString(text)

// Clean up whitespace
text = strings.TrimSpace(text)
re = regexp.MustCompile(`\s+`)
text = re.ReplaceAllString(text, " ")

return text
}

// GenerateRandomString generates a random string of specified length
func (h *Helper) GenerateRandomString(length int) (string, error) {
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
bytes := make([]byte, length)

if _, err := rand.Read(bytes); err != nil {
return "", err
}

for i := range bytes {
bytes[i] = charset[bytes[i]%byte(len(charset))]
}

return string(bytes), nil
}

// GenerateToken generates a secure random token
func (h *Helper) GenerateToken() (string, error) {
return h.GenerateRandomString(32)
}

// ExtractExcerpt extracts an excerpt from content
func (h *Helper) ExtractExcerpt(content string, maxLength int) string {
if maxLength <= 0 {
maxLength = 150
}

// Remove HTML tags first
text := h.SanitizeHTML(content)

// Truncate if necessary
if len(text) <= maxLength {
return text
}

// Find the last space before maxLength to avoid cutting words
excerpt := text[:maxLength]
lastSpace := strings.LastIndex(excerpt, " ")
if lastSpace > 0 {
excerpt = excerpt[:lastSpace]
}

return excerpt + "..."
}

// CalculateReadingTime estimates reading time in minutes
func (h *Helper) CalculateReadingTime(content string) int {
// Average reading speed: 200 words per minute
wordsPerMinute := 200

// Remove HTML and count words
text := h.SanitizeHTML(content)
words := strings.Fields(text)
wordCount := len(words)

readingTime := wordCount / wordsPerMinute
if readingTime < 1 {
readingTime = 1
}

return readingTime
}

// FormatFileSize formats file size in bytes to human readable format
func (h *Helper) FormatFileSize(bytes int64) string {
const unit = 1024
if bytes < unit {
return fmt.Sprintf("%d B", bytes)
}

div, exp := int64(unit), 0
for n := bytes / unit; n >= unit; n /= unit {
div *= unit
exp++
}

return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// MaskEmail masks email address for privacy
func (h *Helper) MaskEmail(email string) string {
parts := strings.Split(email, "@")
if len(parts) != 2 {
return email
}

username := parts[0]
domain := parts[1]

if len(username) <= 2 {
return email
}

masked := username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
return masked + "@" + domain
}

// GenerateSlug creates a URL-friendly slug with uniqueness consideration
func (h *Helper) GenerateSlugWithSuffix(title string, existingSlugs []string) string {
baseSlug := h.createBasicSlug(title)

// Check if slug already exists
slug := baseSlug
counter := 1

for h.containsSlug(existingSlugs, slug) {
slug = fmt.Sprintf("%s-%d", baseSlug, counter)
counter++
}

return slug
}

// createBasicSlug creates a basic slug from title
func (h *Helper) createBasicSlug(title string) string {
slug := strings.ToLower(title)
slug = strings.TrimSpace(slug)

// Replace spaces and special characters with hyphens
re := regexp.MustCompile(`[^a-z0-9\s-]`)
slug = re.ReplaceAllString(slug, "")

re = regexp.MustCompile(`[\s-]+`)
slug = re.ReplaceAllString(slug, "-")

// Remove leading/trailing hyphens
slug = strings.Trim(slug, "-")

return slug
}

// containsSlug checks if slug exists in the slice
func (h *Helper) containsSlug(slugs []string, slug string) bool {
for _, s := range slugs {
if s == slug {
return true
}
}
return false
}

// TimeAgo returns a human-readable representation of time elapsed
func (h *Helper) TimeAgo(t time.Time) string {
duration := time.Since(t)

switch {
case duration < time.Minute:
return "just now"
case duration < time.Hour:
minutes := int(duration.Minutes())
if minutes == 1 {
return "1 minute ago"
}
return fmt.Sprintf("%d minutes ago", minutes)
case duration < 24*time.Hour:
hours := int(duration.Hours())
if hours == 1 {
return "1 hour ago"
}
return fmt.Sprintf("%d hours ago", hours)
case duration < 30*24*time.Hour:
days := int(duration.Hours() / 24)
if days == 1 {
return "1 day ago"
}
return fmt.Sprintf("%d days ago", days)
case duration < 365*24*time.Hour:
months := int(duration.Hours() / (24 * 30))
if months == 1 {
return "1 month ago"
}
return fmt.Sprintf("%d months ago", months)
default:
years := int(duration.Hours() / (24 * 365))
if years == 1 {
return "1 year ago"
}
return fmt.Sprintf("%d years ago", years)
}
}
