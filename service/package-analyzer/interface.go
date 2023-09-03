package packageAnalyzer

// PackageManager defines the interface for fetching package information.
type PackageManager interface {
	Name() string
	Detect(directory string) bool
	FetchPackages() ([]*PackageInfo, error)
	FetchInfo(packageName, packageVersion string) (*PackageInfo, error)
	ParseArgs([]string) (bool, *[]PackageInfo, error)
}

type PackageInfo struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Author      string        `json:"author"`
	Maintainers []string      `json:"maintainers"`
	Funding     []FundingInfo `json:"funding"`
	Description string        `json:"description"`
}
