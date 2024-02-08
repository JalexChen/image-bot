package managers

import (
	"fmt"
	"regexp"
	"strings"
)

func CleanTags(ver string) (string, error) {
	ver = strings.ReplaceAll(ver, "_", ".")
	ver = strings.ReplaceAll(ver, "\"", "")
	regex := regexp.MustCompile(`[a-z]?(\d+\.\d+\.\d+)`)
	match := regex.FindStringSubmatch(ver)
	if match[1] == "" {
		return "", fmt.Errorf("unable to clean tag")
	}
	version := match[1]
	return version, nil
}
