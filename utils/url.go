package utils

import (
	"fmt"
	"strings"
)

func GenerateImageURL(host, publicID string) string {
	scheme := "https"
	if strings.Contains(host, "localhost") {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/image/%s", scheme, host, publicID)
}
