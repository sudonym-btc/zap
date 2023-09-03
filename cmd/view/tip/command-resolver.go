package tipView

import (
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
var PackageCommandResolvers = funk.Map(packageAnalyzer.PackageManagers, func(pm packageAnalyzer.PackageManager) CommandResolver {
	return CommandResolver{
		Model: func(m TipModel) task.ModelI {
			return InitialPackageManagerModel(pm, &m.name, &m.version, m)
		},
		Command:   pm.Name(),
		Detect:    pm.Detect,
		ParseArgs: pm.ParseArgs,
	}
}).([]CommandResolver)

var CommandResolvers = append(SpecialCommandResolvers, PackageCommandResolvers...)

type CommandResolver struct {
	Model     func(TipModel) task.ModelI
	Command   string
	Detect    func(string) bool
	ParseArgs func([]string) (bool, *[]packageAnalyzer.PackageInfo, error)
}

// type PackageCommandResolver struct {
// 	CommandResolver
// 	ParseArgs func([]string) (bool, *[]packageAnalyzer.PackageInfo, error)
// }
