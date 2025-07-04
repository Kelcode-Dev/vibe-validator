package scanner

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "path/filepath"
  "strings"
  "os"
)

// PyDeps maps dependency name to list of file paths where found
type PyDeps map[string][]string

// ScanPython scans the given project directory for Python ecosystem deps
func ScanPython(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (PyDeps, error) {
  if verbosity >= 2 {
    fmt.Println("Scanning Python...")
  }
  deps := make(PyDeps)

  err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return nil
    }

    if info.IsDir() {
      base := filepath.Base(path)
      if !includeVendor && (base == ".venv" || base == "env" || base == "__pycache__") {
        return filepath.SkipDir
      }
      return nil
    }

    base := filepath.Base(path)
    switch base {
    case "requirements.txt":
      parseRequirements(path, deps)
    case "Pipfile.lock":
      if includeLockfiles {
        parsePipfileLock(path, deps)
      }
    }
    return nil
  })

  if err != nil {
    return nil, err
  }

  if verbosity >= 2 {
    fmt.Printf("Finished scanning Python... %d deps found\n\n", len(deps))
  }
  return deps, nil
}

// parseRequirements extracts package names from requirements.txt
func parseRequirements(path string, deps PyDeps) error {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }

  lines := strings.Split(string(data), "\n")
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") {
      continue
    }
    // Strip extras and version operators
    name := line
    if idx := strings.Index(name, "["); idx != -1 {
      name = name[:idx]
    }
    for _, sep := range []string{">=", "==", "<=", "~=", "!=", ">", "<"} {
      if idx := strings.Index(name, sep); idx != -1 {
        name = name[:idx]
      }
    }
    name = strings.TrimSpace(name)
    if name != "" {
      deps[name] = append(deps[name], path)
    }
  }
  return nil
}

// parsePipfileLock extracts package names from Pipfile.lock
func parsePipfileLock(path string, deps PyDeps) error {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }

  var lock struct {
    Default map[string]interface{} `json:"default"`
    Develop map[string]interface{} `json:"develop"`
  }

  if err := json.Unmarshal(data, &lock); err != nil {
    return err
  }

  collectDeps := func(depMap map[string]interface{}) {
    for name := range depMap {
      deps[name] = append(deps[name], path)
    }
  }

  collectDeps(lock.Default)
  collectDeps(lock.Develop)

  return nil
}
