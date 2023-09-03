package packageAnalyzer

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/gookit/slog"
	"github.com/thoas/go-funk"
)

// NpmPackage represents npm package information.
type NpmListPackages struct {
	dependencies NpmPackageDependency `json:"dependencies"`
}
type NpmPackageDependency struct {
}

// NpmPackage represents npm package information.
type NpmPackage struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Author      string        `json:"author"`
	Maintainers []string      `json:"maintainers"`
	Funding     []FundingInfo `json:"funding"`
	Description string        `json:"description"`
}
type NpmPackageSingleFunding struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Maintainers []string    `json:"maintainers"`
	Funding     FundingInfo `json:"funding"`
	Description string      `json:"description"`
}

// FundingInfo represents funding information.
type FundingInfo struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// NpmManager is a struct that implements the PackageManager interface for npm.
type NpmManager struct{}
type NpmCommandResolver struct{}

func (*NpmManager) Name() string {
	return "npm"
}

func (npm *NpmManager) ParseArgs(args []string) (bool, *[]PackageInfo, error) {
	slog.Debug("Parsing args", args)
	if len(args) > 0 && args[0] == "install" {
		slog.Debug("Found npm install command")
		args = args[1:]
		if len(args) > 0 {
			packages := funk.Map(args, func(arg string) PackageInfo {
				return PackageInfo{Name: arg}
			}).([]PackageInfo)
			return true, &packages, nil
		}
		return true, nil, nil
	}
	return false, nil, nil
}

func (npm *NpmManager) Detect(directory string) bool {
	// Check if package.json exists in the specified directory
	cmd := exec.Command("ls", directory)
	output, _ := cmd.Output()
	return strings.Contains(string(output), "package.json")
}

func (npm *NpmManager) FetchPackages() ([]*PackageInfo, error) {
	slog.Debug("Fetching npm packages")
	cmd := exec.Command("npm", "ls", "--json")
	output, err := cmd.Output()

	if err != nil && err.Error() != "exit status 1" {
		slog.Warn("Failed fetching npm packages", output, err)
		return nil, err
	}
	slog.Debug("Fetched npm packages", string(output))

	var z map[string]interface {
	}

	if err := json.Unmarshal(output, &z); err != nil {
		return []*PackageInfo{}, err
	}

	return funk.Map((z["dependencies"]), func(name string, data interface{}) *PackageInfo {
		version := data.(map[string]interface{})["version"]
		if version == nil {
			version = data.(map[string]interface{})["required"]
		}
		return &PackageInfo{Name: name, Version: version.(string)}
	}).([]*PackageInfo), nil
}

func (npm *NpmManager) FetchInfo(packageName, packageVersion string) (*PackageInfo, error) {
	slog.Debug("Fetching npm package info", packageName, packageVersion)
	cmd := exec.Command("npm", "show", packageName, "--json")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("Failed fetching npm package info", output, err)
		return nil, err
	}
	slog.Debug("Fetched npm package info", string(output))

	var packageInfo NpmPackage
	err = json.Unmarshal(output, &packageInfo)
	if err != nil {
		var packageInfoFunding NpmPackageSingleFunding
		err2 := json.Unmarshal(output, &packageInfoFunding)
		if err2 != nil {
			return nil, err
		}
		packageInfo = NpmPackage{
			Name:        packageInfoFunding.Name,
			Version:     packageInfoFunding.Version,
			Maintainers: packageInfoFunding.Maintainers,
			Funding:     []FundingInfo{packageInfoFunding.Funding},
			Description: packageInfoFunding.Description,
			Author:      packageInfoFunding.Author}

	}

	return &PackageInfo{
		Author:      packageInfo.Author,
		Name:        packageInfo.Name,
		Description: packageInfo.Description,
		Version:     packageInfo.Version,
		Maintainers: packageInfo.Maintainers,
		Funding:     packageInfo.Funding,
	}, nil
}
