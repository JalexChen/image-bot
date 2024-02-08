package managers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v56/github"
)

var client = github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

// GetReleaseList can be used if there are multiple versions of something we are looking for
func GetReleaseList(ctx context.Context, owner, repo string, versions map[string]string) (map[string]string, error) {
	releaseMap := make(map[string]string)
	releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("releases not found with error: %v and status code: %d", err, resp.StatusCode)
	}

	for key, versionString := range versions {
		for _, release := range releases {
			releaseVersion := release.GetTagName()
			if strings.Contains(releaseVersion, versionString) {
				version, err := CleanTags(releaseVersion)
				if err != nil {
					return nil, fmt.Errorf("unable to clean tag: %v", err)
				}
				releaseMap[key] = version
				break
			}
		}
	}

	if len(releaseMap) != len(versions) {
		return nil, fmt.Errorf("not all versions found: %v", releaseMap)
	}
	return releaseMap, nil
}

// GetLatestRelease gets the latest release from github's releases endpoint
func GetLatestRelease(ctx context.Context, owner, repo, reference string) (string, error) {
	releases, resp, err := client.Repositories.GetLatestRelease(ctx, owner, repo)

	if resp.StatusCode == http.StatusOK {
		release := releases.GetTagName()
		version, err := CleanTags(release)
		if err != nil {
			return "", fmt.Errorf("unable to clean tag: %v", err)
		}
		// Gradle packages do not follow semver and strip the ".0" patch version for their downloads, resulting in broken links
		if owner == "gradle" && version[len(version)-2:] == ".0" {
			version = strings.TrimSuffix(version, ".0")
		}
		return version, nil
	}

	if err != nil && reference != "" {
		fmt.Println("Release endpoint failure; trying Ref")
		version, err := getTagFromRef(ctx, owner, repo, reference)
		if err != nil {
			return "", fmt.Errorf("failed to get tag from ref: %v", err)
		}
		return version, nil
	}
	return "", fmt.Errorf("the release endpoint was not found")
}

// GetTagFromRef can be used by iteself or as a fallback when repos do not have releases or if you are looking for a specific tag
// For example, gcloud has 4xx versions, so setting the ref in your configuration as "4" would return all instances of a tag that start with 4
// Note: refs are in chronological order, with the most recent entry at the end of the array
func getTagFromRef(ctx context.Context, owner, repo, reference string) (string, error) {
	opt := &github.ReferenceListOptions{Ref: fmt.Sprintf("tags/%s", reference)}

	ref, resp, err := client.Git.ListMatchingRefs(ctx, owner, repo, opt)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ref endpoint not available: %v", err)
	}

	matchedRef := ref[len(ref)-1:][0].Ref
	ver := strings.Split(*matchedRef, "/")
	version, err := CleanTags(ver[2])
	if err != nil {
		return "", fmt.Errorf("unable to clean tag: %v", err)
	}
	return version, nil
}
