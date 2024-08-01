package internal

import (
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "Invalid URL format", err
	}
	validURL := parsedURL.String()
	return validURL, nil
}
