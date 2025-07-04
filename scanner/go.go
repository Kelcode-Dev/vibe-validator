package scanner

import (
  "golang.org/x/mod/modfile"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
)

// GoDeps maps module names to list of file paths where found
type GoDeps map[string][]string

func ScanGo(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (GoDeps, error) {
  if verbosity >= 2 {
    fmt.Println("Scanning Go...")
  }
  deps := make(GoDeps)

  err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return nil
    }

    if info.IsDir() {
      base := filepath.Base(path)
      if !includeVendor && (base == "vendor" || base == "Godeps") {
        return filepath.SkipDir
      }
      return nil
    }

    if filepath.Base(path) == "go.mod" {
      modDeps, err := parseGoMod(path)
      if err == nil {
        for mod := range modDeps {
          deps[mod] = append(deps[mod], path)
        }
      }
    }

    return nil
  })

  if err != nil {
    return nil, err
  }

  if verbosity >= 2 {
    fmt.Printf("Finished scanning Go... %d deps found\n\n", len(deps))
  }
  return deps, nil
}

func parseGoMod(path string) (map[string]struct{}, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  f, err := modfile.Parse(path, data, nil)
  if err != nil {
    return nil, err
  }

  mods := make(map[string]struct{})
  for _, req := range f.Require {
    mods[req.Mod.Path] = struct{}{}
  }
  return mods, nil
}
