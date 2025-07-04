package scanner

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/pelletier/go-toml"
)

// RustDeps maps crate name to list of file paths where found
type RustDeps map[string][]string

func ScanRust(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (RustDeps, error) {
  if verbosity >= 2 {
    fmt.Println("Scanning Ruby...")
  }
  deps := make(RustDeps)

  err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return nil
    }

    if info.IsDir() {
      base := filepath.Base(path)
      if !includeVendor && base == "target" {
        return filepath.SkipDir
      }
      return nil
    }

    base := filepath.Base(path)
    switch base {
    case "Cargo.toml":
      parseCargoToml(path, deps)
    case "Cargo.lock":
      if includeLockfiles {
        parseCargoLock(path, deps)
      }
    }
    return nil
  })

  if err != nil {
    return nil, err
  }

  if verbosity >= 2 {
    fmt.Printf("Finished scanning Rust... %d deps found\n\n", len(deps))
  }
  return deps, nil
}

// parseCargoToml extracts dependencies from Cargo.toml [dependencies] and [dev-dependencies]
func parseCargoToml(path string, deps RustDeps) error {
  content, err := os.ReadFile(path)
  if err != nil {
    return err
  }

  tree, err := toml.LoadBytes(content)
  if err != nil {
    return err
  }

  // Extract [dependencies]
  if depsTable := tree.Get("dependencies"); depsTable != nil {
    if m, ok := depsTable.(*toml.Tree); ok {
      for _, key := range m.Keys() {
        deps[key] = append(deps[key], path)
      }
    }
  }

  // Extract [dev-dependencies]
  if devTable := tree.Get("dev-dependencies"); devTable != nil {
    if m, ok := devTable.(*toml.Tree); ok {
      for _, key := range m.Keys() {
        deps[key] = append(deps[key], path)
      }
    }
  }

  return nil
}

// parseCargoLock extracts dependencies from Cargo.lock file
func parseCargoLock(path string, deps RustDeps) error {
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  var currentPkg string
  inPackage := false

  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "[[package]]" {
      inPackage = true
      currentPkg = ""
      continue
    }
    if inPackage {
      if strings.HasPrefix(line, "name = ") {
        currentPkg = strings.Trim(line[len("name = "):], "\"")
      }
      if line == "" {
        if currentPkg != "" {
          deps[currentPkg] = append(deps[currentPkg], path)
        }
        inPackage = false
      }
    }
  }
  // Handle last package if file doesn’t end with blank line
  if inPackage && currentPkg != "" {
    deps[currentPkg] = append(deps[currentPkg], path)
  }

  return scanner.Err()
}
