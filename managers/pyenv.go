package managers

import (
	"context"
	"fmt"
	"os/exec"
)

func RunPyenv(ctx context.Context, version string) (map[string]string, error) {
	pythonMap := make(map[string]string)

	cmd := exec.Command("pyenv", "latest", "--known", version)
	pythonVersion, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("something went wrong with the pyenv command: %v", err)
	}

	pythonMap[fmt.Sprintf("python%s", version)] = string(pythonVersion)
	return pythonMap, nil
}
