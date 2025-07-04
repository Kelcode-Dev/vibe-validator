package scanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// JsDeps maps dependency name to list of file paths where found
type JsDeps map[string][]string

// ScanJavaScript scans the given project directory for JS ecosystem deps
func ScanJavaScript(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (JsDeps, error) {
	if verbosity >= 2 {
		fmt.Println("Scanning JavaScript...")
	}
	deps := make(JsDeps)

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			base := filepath.Base(path)
			if !includeVendor && base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		base := filepath.Base(path)
		switch base {
		case "package.json":
			parsePackageJSON(path, deps)

		case "package-lock.json", "yarn.lock", "pnpm-lock.yaml":
			if includeLockfiles {
				parseLockfile(path, deps)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if verbosity >= 2 {
		fmt.Printf("Finished scanning JavaScript... %d deps found\n\n", len(deps))
	}
	return deps, nil
}

// parsePackageJSON extracts dependencies & devDependencies from package.json
func parsePackageJSON(path string, deps JsDeps) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var obj struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	for name := range obj.Dependencies {
		deps[name] = append(deps[name], path)
	}
	for name := range obj.DevDependencies {
		deps[name] = append(deps[name], path)
	}

	return nil
}

// parseLockfile parses npm/yarn/pnpm lockfiles to gather all locked deps
func parseLockfile(path string, deps JsDeps) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// npm lockfile structure
	var lock struct {
		Dependencies map[string]interface{} `json:"dependencies"`
	}

	if err := json.Unmarshal(data, &lock); err != nil {
		return err
	}

	// recursively collect dependencies
	var collect func(map[string]interface{})
	collect = func(depsMap map[string]interface{}) {
		for name, val := range depsMap {
			deps[name] = append(deps[name], path)
			if depInfo, ok := val.(map[string]interface{}); ok {
				if nestedDeps, ok := depInfo["dependencies"].(map[string]interface{}); ok {
					collect(nestedDeps)
				}
			}
		}
	}
	collect(lock.Dependencies)

	return nil
}
