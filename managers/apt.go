package managers

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GetAptVersions retrieves the latest versions from the apt repository and returns the latest version of a given package
func GetAptVersions(ctx context.Context, packages []string) (map[string]string, error) {
	aptMap := make(map[string]string)
	updateCmd := exec.Command("sudo", "apt-get", "update")
	if err := updateCmd.Run(); err != nil {
		return nil, fmt.Errorf("error running apt-get update: %v", err)
	}

	// apt-cache policy gets the latest package versions within the apt repository.
	// this will occassionally be mismatched with the most current release on github, so we use apt's repo to reduce the risk of breaking changes
	for _, pkg := range packages {
		cmd := exec.Command("sudo", "apt-cache", "policy", pkg)
		versions, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("error running apt-cache policy: %v", err)
		}

		output := string(versions)
		version, err := extractLatestVersion(output)
		if err != nil {
			return nil, fmt.Errorf("error getting apt-cache policy: %v", err)
		}
		aptMap[pkg] = version
	}

	if len(aptMap) != len(packages) {
		return nil, fmt.Errorf("aptMap does not have the correct number of packages: %v", aptMap)
	}
	return aptMap, nil
}

func extractLatestVersion(output string) (string, error) {
	re := regexp.MustCompile(`Version table:\s+([^\n]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) != 2 {
		return "", fmt.Errorf("no match found")
	}
	match := strings.Split(matches[1], " ")

	return match[0], nil
}
