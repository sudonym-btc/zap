package packageAnalyzer

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os/exec"
// 	"strings"

// 	"github.com/thoas/go-funk"
// )

// // PipPackage represents pip package information.
// type PipPackage struct {
// 	Name        string `json:"name"`
// 	Version     string `json:"version"`
// 	Author      string `json:"author"`
// 	Description string `json:"description"`
// }

// // PipManager is a struct that implements the PackageManager interface for pip.
// type PipManager struct{}

// func (*PipManager) Name() string {
// 	return "pip"
// }

// func (*PipManager) ParseArgs(args []string) (bool, *[]PackageInfo, error) {
// 	if args[0] == "install" {
// 		args = args[1:]
// 		if len(args) > 0 {
// 			return true, funk.Map(args, func(arg string) PackageInfo {
// 				return PackageInfo{Name: arg}
// 			}).(*[]PackageInfo), nil
// 		}
// 		return true, nil, nil
// 	}
// 	return false, nil, nil
// }

// func (pip *PipManager) Detect(directory string) bool {
// 	// Check if requirements.txt exists in the specified directory
// 	cmd := exec.Command("ls", directory)
// 	output, _ := cmd.Output()
// 	return contains(output, "requirements.txt")
// }

// func (pip *PipManager) FetchPackages() ([]*PackageInfo, error) {
// 	cmd := exec.Command("npm", "ls", "--json")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var z map[string]interface{}

// 	if err := json.Unmarshal(output, &z); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(z)

// 	return funk.Map(z["dependencies"], func(x map[string]interface{}) *PackageInfo {
// 		r, _ := pip.FetchInfo(x["name"].(string), x["version"].(string))
// 		return r
// 	}).([]*PackageInfo), nil
// }

// func (pip *PipManager) FetchInfo(packageName, packageVersion string) (*PackageInfo, error) {
// 	cmd := exec.Command("pip", "show", packageName)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	packageInfo := parsePipOutput(output)
// 	return &PackageInfo{
// 		Name:    packageInfo.Name,
// 		Version: packageInfo.Version,
// 	}, nil
// }

// func parsePipOutput(output []byte) *PipPackage {
// 	lines := strings.Split(string(output), "\n")
// 	pkgInfo := PipPackage{}
// 	for _, line := range lines {
// 		parts := strings.SplitN(line, ":", 2)
// 		if len(parts) == 2 {
// 			key := strings.TrimSpace(parts[0])
// 			value := strings.TrimSpace(parts[1])
// 			switch key {
// 			case "Name":
// 				pkgInfo.Name = value
// 			case "Version":
// 				pkgInfo.Version = value
// 			case "Author":
// 				pkgInfo.Author = value
// 			case "Summary":
// 				pkgInfo.Description = value
// 			}
// 		}
// 	}
// 	return &pkgInfo
// }

// func contains(list []byte, target string) bool {
// 	return strings.Contains(string(list), target)
// }
