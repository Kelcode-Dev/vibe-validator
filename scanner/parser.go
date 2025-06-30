package scanner

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"

  "golang.org/x/mod/modfile"
)

// ScanDependencies returns three maps of dependency name → file path
func ScanDependencies(projectPath string, includeLockfiles bool, includeVendor bool) (map[string][]string, map[string][]string, map[string][]string) {
  pyDeps := make(map[string][]string)
  npmDeps := make(map[string][]string)
  goDeps := make(map[string][]string)

  excludedDirs := map[string]struct{}{
    "node_modules": {},
    "vendor":       {},
    ".venv":        {},
    "env":          {},
    "__pycache__":  {},
  }

  filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return nil
    }

    if info.IsDir() {
      base := filepath.Base(path)
      if !includeVendor {
        if _, excluded := excludedDirs[base]; excluded {
          return filepath.SkipDir
        }
      }
      return nil
    }

    switch filepath.Base(path) {
      case "requirements.txt":
        parseRequirements(path, pyDeps)

      case "Pipfile.lock":
        if includeLockfiles {
          parsePipfileLock(path, pyDeps)
        }

      case "package.json":
        parsePackageJSON(path, npmDeps)

      case "package-lock.json", "yarn.lock", "pnpm-lock.yaml":
        if includeLockfiles {
          parseNpmLockfile(path, npmDeps)
        }

      case "go.mod":
        if mods, err := parseGoMod(path); err == nil {
          for mod := range mods {
            goDeps[mod] = append(goDeps[mod], path)
          }
        }
    }

    return nil
  })

  return pyDeps, npmDeps, goDeps
}

// parseRequirements reads a requirements.txt, strips versions, and populates deps[name]=path
func parseRequirements(path string, deps map[string][]string) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return
  }

  for _, line := range strings.Split(string(data), "\n") {
    line = strings.TrimSpace(line)

    if line == "" || strings.HasPrefix(line, "#") {
      continue
    }

    // strip version operators
    name := line
    for _, sep := range []string{">=", "==", "<=", "~=", "!=", ">", "<", "["} {
      if idx := strings.Index(name, sep); idx != -1 {
        name = name[:idx]
      }
    }
    name = strings.TrimSpace(name)

    if name != "" {
      deps[name] = append(deps[name], path)
    }
  }
}

// parsePackageJSON reads package.json and populates deps[name]=path
func parsePackageJSON(path string, deps map[string][]string) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return
  }

  var obj struct {
    Dependencies    map[string]interface{} `json:"dependencies"`
    DevDependencies map[string]interface{} `json:"devDependencies"`
  }

  if err := json.Unmarshal(data, &obj); err != nil {
    return
  }

  for name := range obj.Dependencies {
    deps[name] = append(deps[name], path)
  }

  for name := range obj.DevDependencies {
    deps[name] = append(deps[name], path)
  }
}

// parseGoMod uses golang.org/x/mod/modfile to extract only 'require' entries
func parseGoMod(path string) (map[string][]string, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  f, err := modfile.Parse(path, data, nil)
  if err != nil {
    return nil, err
  }
  mods := make(map[string][]string)
  for _, req := range f.Require {
    mods[req.Mod.Path] = append(mods[req.Mod.Path], path)
  }
  return mods, nil
}

// parsePipLockfile checks yarn.lock, package-lock.json and pnpm-lock.yaml for dependencies
func parsePipfileLock(path string, deps map[string][]string) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return
  }

  var lockfile struct {
    Default map[string]interface{} `json:"default"`
    Develop map[string]interface{} `json:"develop"`
  }

  if err := json.Unmarshal(data, &lockfile); err != nil {
    return
  }

  // Add all default packages
  for pkg := range lockfile.Default {
    deps[pkg] = append(deps[pkg], path)
  }
  // Add all dev packages
  for pkg := range lockfile.Develop {
    deps[pkg] = append(deps[pkg], path)
  }
}

// parseNpmLockfile checks yarn.lock, package-lock.json and pnpm-lock.yaml for dependencies
func parseNpmLockfile(path string, deps map[string][]string) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return
  }

  var lock struct {
    Dependencies map[string]interface{} `json:"dependencies"`
  }
  if err := json.Unmarshal(data, &lock); err != nil {
    return
  }

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
}
