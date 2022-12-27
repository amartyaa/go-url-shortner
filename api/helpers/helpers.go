package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Params struct {
	URL   string
	Check bool
}

func EnforeceHTTPS(url string) string {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func LoopDomain(url string) bool {
	if url == os.Getenv("APP_HOST") {
		fmt.Println("Same domain")
		return false
	}
	newURL := strings.Replace(url, "https://", "", 1)
	newURL = strings.Replace(newURL, "http://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	fmt.Println(newURL)
	if newURL == os.Getenv("APP_HOST") || strings.Contains(newURL, os.Getenv("APP_HOST")) {
		fmt.Println("Loop domain")
		return false
	}
	return true
}

func GenerateID() string {
	return uuid.New().String()[0:6]
}
func ValidParams(url string) Params {
	var params Params
	params.URL = url
	params.Check = false
	// if strings.Contains(url, os.Getenv("APP_HOST")) {
	// 	url = strings.Split(url, os.Getenv("APP_HOST"))[1]
	// } 
	if strings.Contains(url, "=") {
		url = strings.Split(url, "=")[1]
		if strings.Contains(url,"/")  {
			url = strings.Split(string(url), "/")[1]
		}
		params.Check = true
		params.URL = url
	}

	return params
}
