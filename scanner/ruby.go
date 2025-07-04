package scanner

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"
)

// RubyDeps maps gem name to list of file paths where found
type RubyDeps map[string][]string

func ScanRuby(projectPath string, includeLockfiles, includeVendor bool, verbosity int) (RubyDeps, error) {
  if verbosity >= 2 {
    fmt.Println("Scanning Ruby...")
  }
  deps := make(RubyDeps)

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
    case "Gemfile":
      parseGemfile(path, deps)
    case "Gemfile.lock":
      if includeLockfiles {
        parseGemfileLock(path, deps)
      }
    }
    return nil
  })

  if err != nil {
    return nil, err
  }

  if verbosity >= 2 {
    fmt.Printf("Finished scanning Ruby... %d deps found\n\n", len(deps))
  }
  return deps, nil
}

// parseGemfile extracts gem names from Gemfile by looking for lines like: gem 'name', 'version'
func parseGemfile(path string, deps RubyDeps) error {
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if strings.HasPrefix(line, "gem ") {
      // crude parsing: extract first string literal after gem keyword
      parts := strings.Split(line, "'")
      if len(parts) >= 2 {
        gemName := parts[1]
        deps[gemName] = append(deps[gemName], path)
      }
    }
  }

  return scanner.Err()
}

// parseGemfileLock extracts gems from Gemfile.lock by parsing lines under "GEM" section
func parseGemfileLock(path string, deps RubyDeps) error {
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  inGemSection := false

  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "GEM" {
      inGemSection = true
      continue
    }
    if inGemSection {
      if line == "" {
        break
      }
      // line format:   gem_name (version)
      if strings.HasPrefix(line, "  ") && strings.Contains(line, " ") {
        gemName := strings.Fields(line)[0]
        deps[gemName] = append(deps[gemName], path)
      }
    }
  }

  return scanner.Err()
}
