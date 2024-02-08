package managers

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestCleanTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "maven",
			input: "maven-3.9.6",
			want:  "3.9.6",
		},
		{
			name:  "remove prefixes",
			input: "random-1.2.3",
			want:  "1.2.3",
		},
		{
			name:  "go test",
			input: "go1.12.4",
			want:  "1.12.4",
		},
		{
			name:  "only semver",
			input: "go1.12.4.12",
			want:  "1.12.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver, err := CleanTags(tt.input)
			if err != nil {
				t.Errorf("got nil value")
			}
			if ver != tt.want {
				t.Errorf("got %s, want %s", ver, tt.want)
			}
		})
	}
}

func TestExtractLatestVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "test matching string",
			input: `Version table:
			5:24.0.7-1~ubuntu.22.04~jammy 500`,
			want: "5:24.0.7-1~ubuntu.22.04~jammy",
		},
		{
			name: "test random string",
			input: `Version table:
			random `,
			want: "random",
		},
		{
			name: "test no new line string",
			input: `Version table:
			`,
			want: "",
		},
		{
			name:  "panic test",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					errMsg := r.(string)
					assert.Equal(t, tt.want, errMsg)
				}
			}()
			extractedValue, err := extractLatestVersion(tt.input)
			if err != nil {
				assert.Equal(t, tt.want, extractedValue)
			}
		})
	}
}

func TestGo(t *testing.T) {
	url := "https://go.dev/dl/?mode=json"

	goMap, err := GetGoVersion(ctx, url)

	if err != nil {
		t.Errorf("GetGoVersion error: %v", err)
	}

	re := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	version := re.FindString(goMap["go"])
	assert.Equal(t, version, goMap["go"])
}

func TestPyenv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "test python 3 major minor",
			input: "3.11",
			want:  "3.11",
		},
		{
			name:  "test python 3 major",
			input: "3",
			want:  "3.",
		},
		{
			name:  "test python 2",
			input: "2",
			want:  "2.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pythonMap, err := RunPyenv(ctx, tt.input)
			if err != nil {
				t.Errorf("get python version errored: %v", err)
			}
			assert.Contains(t, tt.want, pythonMap[tt.input])
		})
	}
}

func TestReleaseList(t *testing.T) {
	tests := []struct {
		name     string
		owner    string
		repo     string
		versions map[string]string
		want     map[string]string
	}{
		{
			name:     "release list node 20, 21",
			owner:    "nodejs",
			repo:     "node",
			versions: map[string]string{"nodelts": "v20", "nodeCurrent": "v21"},
			want:     map[string]string{"nodelts": "20", "nodeCurrent": "21"},
		},
		{
			name:     "release list node 18, 20",
			owner:    "nodejs",
			repo:     "node",
			versions: map[string]string{"nodelts": "v18", "nodeCurrent": "v20"},
			want:     map[string]string{"nodelts": "18", "nodeCurrent": "20"},
		},
	}

	for _, tt := range tests {
		for key := range tt.versions {
			t.Run(tt.name, func(t *testing.T) {
				releaseMap, err := GetReleaseList(ctx, tt.owner, tt.repo, tt.versions)
				if err != nil {
					t.Errorf("get release list errored: %v", err)
				}
				assert.Contains(t, releaseMap[key], tt.want[key])
			})
		}
	}
}

func TestLatestRelease(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		repo  string
		ref   string
		want  string
	}{
		{
			name:  "release list node 20, 21",
			owner: "ruby",
			repo:  "ruby",
			ref:   "",
			want:  "3.",
		},
		{
			name:  "release list node 18, 20",
			owner: "mikefarah",
			repo:  "yq",
			ref:   "",
			want:  "4.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release, err := GetLatestRelease(ctx, tt.owner, tt.repo, tt.ref)
			if err != nil {
				t.Errorf("get latest release errored: %v", err)
			}
			assert.Contains(t, release, tt.want)
		})
	}
}

func TestTagFromRef(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		repo  string
		ref   string
		want  string
	}{
		{
			name:  "GetRef for Node",
			owner: "nodejs",
			repo:  "node",
			ref:   "v21",
			want:  "21.",
		},
		{
			name:  "GetRef for gcloud",
			owner: "twistedpair",
			repo:  "google-cloud-sdk",
			ref:   "4",
			want:  "4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release, err := getTagFromRef(ctx, tt.owner, tt.repo, tt.ref)
			if err != nil {
				t.Errorf("get tag from ref errored: %v", err)
			}
			assert.Contains(t, release, tt.want)
		})
	}
}
