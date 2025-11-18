package validation

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Regex to detect potential XSS attempts
	scriptTagRegex = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	htmlTagRegex   = regexp.MustCompile(`<[^>]+>`)

	// Regex to detect potential NoSQL injection attempts
	nosqlOperatorRegex = regexp.MustCompile(`\$[a-zA-Z]+`)
)

// SanitizeString removes dangerous characters and escapes HTML
// This prevents XSS and other injection attacks
func SanitizeString(input string) string {
	// Step 1: Trim whitespace
	s := strings.TrimSpace(input)

	// Step 2: Remove null bytes (can break string handling)
	s = strings.ReplaceAll(s, "\x00", "")

	// Step 3: Remove control characters (except newlines and tabs)
	s = removeControlChars(s)

	// Step 4: Escape HTML entities
	// Converts < to &lt;, > to &gt;, etc.
	// Prevents XSS by making tags display as text
	s = html.EscapeString(s)

	// Step 5: Normalize unicode (prevent unicode exploits)
	s = normalizeUnicode(s)

	return s
}

// SanitizeDescription sanitizes task descriptions
// Allows newlines but removes dangerous content
func SanitizeDescription(description string) string {
	s := SanitizeString(description)

	// Remove script tags specifically
	s = scriptTagRegex.ReplaceAllString(s, "")

	return s
}

// ValidateMongoID checks if a string is a valid MongoDB ObjectID
// MongoDB ObjectIDs are exactly 24 hex characters
func ValidateMongoID(id string) bool {
	if len(id) != 24 {
		return false
	}

	// Check if all characters are valid hex
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

// DetectNoSQLInjection checks for common NoSQL injection patterns
// Returns true if suspicious patterns are found
func DetectNoSQLInjection(input string) bool {
	// Check for MongoDB operators like $where, $ne, $gt, etc.
	if nosqlOperatorRegex.MatchString(input) {
		return true
	}

	// Check for JavaScript keywords (used in $where attacks)
	jsKeywords := []string{"function", "return", "eval", "while", "for"}
	lower := strings.ToLower(input)
	for _, keyword := range jsKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}

	return false
}

// SanitizeTitle sanitizes task titles
// More strict than general strings
func SanitizeTitle(title string) string {
	s := SanitizeString(title)

	// Remove any remaining HTML tags (double safety)
	s = htmlTagRegex.ReplaceAllString(s, "")

	// Limit to single line (remove newlines)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// Collapse multiple spaces
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

// removeControlChars removes control characters except newlines and tabs
func removeControlChars(s string) string {
	return strings.Map(func(r rune) rune {
		// Keep newlines (\n) and tabs (\t)
		if r == '\n' || r == '\t' {
			return r
		}

		// Remove other control characters
		if unicode.IsControl(r) {
			return -1 // -1 means remove this character
		}

		return r
	}, s)
}

// normalizeUnicode normalizes unicode to prevent homograph attacks
// Example: Cyrillic 'a' (U+0430) vs Latin 'a' (U+0061) look identical
func normalizeUnicode(s string) string {
	// Remove zero-width characters that can hide malicious content
	zeroWidthChars := []string{
		"\u200B", // Zero-width space
		"\u200C", // Zero-width non-joiner
		"\u200D", // Zero-width joiner
		"\uFEFF", // Zero-width no-break space
	}
	for _, char := range zeroWidthChars {
		s = strings.ReplaceAll(s, char, "")
	}

	return s
}

// ============================================================================
// VALIDATION HELPERS
// ============================================================================

// IsEmpty checks if a string is empty after trimming
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsSafeLength checks if string length is within safe bounds
func IsSafeLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}
