package scanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// PhpDeps maps package name to list of file paths where found
type PhpDeps map[string][]string

// ScanPHP scans for PHP Composer dependencies in a project directory
func ScanPHP(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (PhpDeps, error) {
	if verbosity >= 2 {
		fmt.Println("Scanning PHP...")
	}
	deps := make(PhpDeps)

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			base := filepath.Base(path)
			if !includeVendor && base == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		base := filepath.Base(path)
		switch base {
		case "composer.json":
			parseComposerJSON(path, deps)
		case "composer.lock":
			if includeLockfiles {
				parseComposerLock(path, deps)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if verbosity >= 2 {
		fmt.Printf("Finished scanning PHP... %d deps found\n\n", len(deps))
	}
	return deps, nil
}

// parseComposerJSON extracts dependencies from composer.json (require & require-dev)
func parseComposerJSON(path string, deps PhpDeps) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var obj struct {
		Require    map[string]interface{} `json:"require"`
		RequireDev map[string]interface{} `json:"require-dev"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	for name := range obj.Require {
		deps[name] = append(deps[name], path)
	}
	for name := range obj.RequireDev {
		deps[name] = append(deps[name], path)
	}

	return nil
}

// parseComposerLock extracts dependencies from composer.lock
func parseComposerLock(path string, deps PhpDeps) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var lock struct {
		Packages []struct {
			Name string `json:"name"`
		} `json:"packages"`
		DevPackages []struct {
			Name string `json:"name"`
		} `json:"packages-dev"`
	}

	if err := json.Unmarshal(data, &lock); err != nil {
		return err
	}

	for _, pkg := range lock.Packages {
		deps[pkg.Name] = append(deps[pkg.Name], path)
	}
	for _, pkg := range lock.DevPackages {
		deps[pkg.Name] = append(deps[pkg.Name], path)
	}

	return nil
}
