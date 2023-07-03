package helpers

import (
	"regexp"
	"strings"
)

func CreateSlug(s string) string {
	// Convert string to lowercase
	str := strings.ToLower(s)

	// Replace non-alphanumeric characters with a hyphen
	reg := regexp.MustCompile("[^a-z0-9]+")
	str = reg.ReplaceAllString(str, "-")

	// Remove leading and trailing hyphens
	str = strings.Trim(str, "-")

	return str
}
