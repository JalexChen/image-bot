package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Go's github repository does not keep any official recent releases, so the only official
// way to grab this information from their sources is to use their go.dev link
func GetGoVersion(ctx context.Context, url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get requested url at %s: %v", url, err)
	}

	goMap := make(map[string]string)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return goMap, fmt.Errorf("something went wrong with reading the url: %v", err)
	}

	var response []map[string]interface{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return goMap, fmt.Errorf("something went wrong with trying to parse the url: %v", err)
	}

	if len(response) > 0 {
		firstEntry := response[0]
		version := firstEntry["version"].(string)
		version, err = CleanTags(version)
		if err != nil {
			return goMap, fmt.Errorf("unable to clean tag: %v", err)
		}
		goMap["go"] = version
	}
	return goMap, nil
}
