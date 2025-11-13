package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func SanitizeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	return rawURL
}

func IsValidShortCode(code string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9-]{3,20}$`)
	return regex.MatchString(code)
}

func FormatError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
