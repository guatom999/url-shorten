package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// IsValidURL validates if a string is a valid URL
func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// SanitizeURL cleans and normalizes a URL
func SanitizeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	return rawURL
}

// IsValidShortCode validates a short code format
func IsValidShortCode(code string) bool {
	// Allow alphanumeric characters and hyphens, 3-20 characters
	regex := regexp.MustCompile(`^[a-zA-Z0-9-]{3,20}$`)
	return regex.MatchString(code)
}

// FormatError wraps an error with context
func FormatError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
