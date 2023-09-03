package packageAnalyzer

import (
	"encoding/json"
	"os/exec"

	"strings"

	"github.com/gookit/slog"
	"github.com/thoas/go-funk"
)

// NpmPackage represents npm package information.
type YarnListPackages struct {
	dependencies NpmPackageDependency `json:"dependencies"`
}
type YarnPackageDependency struct {
}

// NpmPackage represents npm package information.
type YarnPackage struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Maintainers []string    `json:"maintainers"`
	Funding     FundingInfo `json:"funding"`
	Description string      `json:"description"`
}

// NpmManager is a struct that implements the PackageManager interface for npm.
type YarnManager struct{}

func (*YarnManager) Name() string {
	return "npm"
}

func (npm *YarnManager) Detect(directory string) bool {
	// Check if package.json exists in the specified directory
	cmd := exec.Command("ls", directory)
	output, _ := cmd.Output()
	return strings.Contains(string(output), "package.json")
}

func (npm *YarnManager) FetchPackages() ([]*PackageInfo, error) {
	slog.Debug("Fetching yarn packages")
	cmd := exec.Command("npm", "ls", "--json")
	output, err := cmd.Output()

	if err != nil && err.Error() != "exit status 1" {
		slog.Warn("Failed fetching yarn packages", err)
		return nil, err
	}

	var z map[string]interface {
	}

	if err := json.Unmarshal(output, &z); err != nil {
		slog.Warn("Failed fetching yarn packages", err)
		return nil, err
	}
	slog.Debug("Fetched yarn packages", string(output))

	return funk.Map((z["dependencies"]), func(name string, data interface{}) *PackageInfo {
		version := data.(map[string]interface{})["version"]
		if version == nil {
			version = data.(map[string]interface{})["required"]
		}
		r, _ := npm.FetchInfo(name, version.(string))
		return r
	}).([]*PackageInfo), nil
}

func (npm *YarnManager) FetchInfo(packageName, packageVersion string) (*PackageInfo, error) {
	slog.Debug("Fetching yarn package info", packageName, packageVersion)
	cmd := exec.Command("npm", "show", packageName, "--json")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("Failed fetching yarn package info", err)
		return nil, err
	}

	var packageInfo NpmPackage
	err = json.Unmarshal(output, &packageInfo)
	if err != nil {
		slog.Warn("Failed fetching yarn package info", err)
		return nil, err
	}
	slog.Debug("Fetched yarn package info", string(output))

	return &PackageInfo{
		Author:      packageInfo.Author,
		Name:        packageInfo.Name,
		Maintainers: packageInfo.Maintainers,
	}, nil
}
