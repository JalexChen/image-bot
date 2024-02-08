package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CircleCI-Public/image-bot/managers"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx := context.Background()
	configuration := flag.String("c", "", "input config file here")
	manifest := flag.String("m", "", "input manifest file here")
	flag.Parse()
	config, err := loadConfig(*configuration)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	releaseMap := make(map[string]string)
	for _, gh := range config.Github {
		if len(gh.Versions) != 0 {
			releaseInfo, err := managers.GetReleaseList(ctx, gh.Owner, gh.Repo, gh.Versions)
			if err != nil {
				log.Fatalf("get release list fatal error: %v", err)
			}

			for pkg, version := range releaseInfo {
				releaseMap[pkg] = version
			}
		} else {
			version, err := managers.GetLatestRelease(ctx, gh.Owner, gh.Repo, gh.Ref)
			if err != nil {
				log.Fatalf("unable to get release from ref: %v", err)
			}
			releaseMap[gh.Repo] = version
		}
	}

	pythonInfo, err := managers.RunPyenv(ctx, "3")
	if err != nil {
		log.Fatalf("error running pyenv: %v", err)
	}

	for pkg, version := range pythonInfo {
		releaseMap[pkg] = version
	}

	goInfo, err := managers.GetGoVersion(ctx, config.Go)
	if err != nil {
		log.Fatalf("error getting go version: %v", err)
	}

	for pkg, version := range goInfo {
		releaseMap[pkg] = version
	}

	if len(config.Apt) != 0 {
		aptInfo, err := managers.GetAptVersions(ctx, config.Apt)
		if err != nil {
			log.Fatalf("error with getting apt versions: %v", err)
		}

		for pkg, version := range aptInfo {
			releaseMap[pkg] = version
		}
	}

	err = updateManifest(*manifest, releaseMap)
	if err != nil {
		log.Fatalf("error updating manifest: %v", err)
	}
}

type Config struct {
	Go     string                `yaml:"go"`
	Github map[string]GitHubRepo `yaml:"github"`
	Apt    []string              `yaml:"apt"`
}

type GitHubRepo struct {
	Owner    string            `yaml:"owner"`
	Repo     string            `yaml:"repo"`
	Ref      string            `yaml:"ref,omitempty"`
	Versions map[string]string `yaml:"versions,omitempty"`
}

func loadConfig(path string) (config Config, err error) {
	manifest, err := os.ReadFile(path)

	if err != nil {
		return Config{}, fmt.Errorf("error loading config: %v", err)
	}

	err = yaml.Unmarshal(manifest, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling manifest: %v", err)
	}
	return config, nil
}

func updateManifest(filePath string, releaseMap map[string]string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	var manifest map[string]interface{}
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return fmt.Errorf("error unmarshalling manifest: %v", err)
	}

	for pkg, version := range releaseMap {
		// Ansible does not accept hyphens or dots, which means these must be replaced with another value, like underscores
		if strings.Contains(pkg, ".") || strings.Contains(pkg, "-") {
			pkg = strings.ReplaceAll(pkg, ".", "_")
			pkg = strings.ReplaceAll(pkg, "-", "_")
		}
		manifest[pkg] = strings.TrimSuffix(version, "\n")
	}

	updatedData, err := yaml.Marshal(&manifest)
	if err != nil {
		return fmt.Errorf("error Marshalling updatedData: %v", err)
	}

	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", filePath, err)
	}
	return nil
}
