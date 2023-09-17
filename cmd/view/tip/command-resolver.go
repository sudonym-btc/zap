package tipView

import (
	"fmt"
	"regexp"

	"github.com/gookit/slog"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	packageAnalyzer "github.com/sudonym-btc/zap/service/package-analyzer"
	"github.com/thoas/go-funk"
)

var SpecialCommandResolvers = []CommandResolver{
	{
		Model: func(m TipModel) task.ModelI {
			return initialGithubUserModel(m.name, true, m)
		},
		Command: "github-user",
	}, {
		Model: func(m TipModel) task.ModelI {
			return initialGithubOrganizationModel(&m.name, m)
		},
		Command: "github-org",
	},
	{
		Model: func(m TipModel) task.ModelI {
			return initialGithubRepoModel(&m.name, m)
		},
		Command: "github-repo",
	},
	{
		Model: func(m TipModel) task.ModelI {
			return initialMaintainerModel(m.name, true, m)
		},
		Command: "address",
	}}
var PackageCommandResolvers = append(funk.Map(packageAnalyzer.PackageManagers, func(pm packageAnalyzer.PackageManager) CommandResolver {
	return CommandResolver{
		Model: func(m TipModel) task.ModelI {
			return InitialPackageManagerModel(pm, &m.name, &m.version, m)
		},
		Command:   pm.Name(),
		Detect:    pm.Detect,
		ParseArgs: pm.ParseArgs,
	}
}).([]CommandResolver), CommandResolver{
	Model: func(m TipModel) task.ModelI {
		return initialGithubRepoModel(&m.name, m)
	},
	Command: "git",
	Detect: func(dir string) bool {
		return false
	},
	ParseArgs: func(args []string) (bool, *[]packageAnalyzer.PackageInfo, error) {
		slog.Debug("Parsing args", args)
		if len(args) > 0 && args[0] == "clone" {
			slog.Debug("Found git clone command")
			args = args[1:]
			if len(args) > 0 {
				packages := funk.Map(args, func(arg string) packageAnalyzer.PackageInfo {
					return packageAnalyzer.PackageInfo{Name: extractOrgRepo(arg)}
				}).([]packageAnalyzer.PackageInfo)
				return true, &packages, nil
			}
			return true, nil, nil
		}
		return false, nil, nil
	},
}, CommandResolver{
	Model: func(m TipModel) task.ModelI {
		return initialGithubRepoModel(&m.name, m)
	},
	Command: "gh",
	Detect: func(dir string) bool {
		return false
	},
	ParseArgs: func(args []string) (bool, *[]packageAnalyzer.PackageInfo, error) {
		slog.Debug("Parsing args", args)
		if len(args) > 0 && args[0] == "repo" && args[1] == "clone" {
			slog.Debug("Found git clone command")
			args = args[2:]
			if len(args) > 0 {
				packages := funk.Map(args, func(arg string) packageAnalyzer.PackageInfo {
					return packageAnalyzer.PackageInfo{Name: arg}
				}).([]packageAnalyzer.PackageInfo)
				return true, &packages, nil
			}
			return true, nil, nil
		}
		return false, nil, nil
	},
})

var CommandResolvers = append(SpecialCommandResolvers, PackageCommandResolvers...)

type CommandResolver struct {
	Model     func(TipModel) task.ModelI
	Command   string
	Detect    func(string) bool
	ParseArgs func([]string) (bool, *[]packageAnalyzer.PackageInfo, error)
}

func extractOrgRepo(url string) string {
	// Define the regular expression pattern to match GitHub URLs
	pattern := `(?:https:\/\/github\.com\/|git@github\.com:)([\w-]+)\/([\w-]+)\.git`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Find the matches in the URL
	matches := regex.FindStringSubmatch(url)

	// If there's a match, return the org/repo part, otherwise return an empty string
	if len(matches) == 3 {
		org := matches[1]
		repo := matches[2]
		return fmt.Sprintf("%s/%s", org, repo)
	}

	return ""
}
