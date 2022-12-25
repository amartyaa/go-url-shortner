package helpers

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func EnforeceHTTPS(url string) string {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func LoopDomain(url string) bool {
	if url == os.Getenv("APP_DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "https://", "", 1)
	newURL = strings.Replace(newURL, "http://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	if newURL == os.Getenv("APP_DOMAIN") || strings.Contains(newURL, os.Getenv("APP_DOMAIN")) {
		return false
	}
	return true
}

func GenerateID() string {
	return uuid.New().String()[0:6]
}
