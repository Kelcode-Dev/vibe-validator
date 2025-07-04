package validator

// ValidationResult holds dependency check results with multiple paths
type ValidationResult struct {
  Name    string
  Source  string   // "npm", "pypi", "go"
  Status  string   // "safe", "investigate", "not_found"
  Details string
  Paths   []string `json:"paths"`
}

func ValidatePackages(allDeps map[string]map[string][]string) []ValidationResult {
  var results []ValidationResult

  for eco, deps := range allDeps {
    for pkg, paths := range deps {
      switch eco {
        case "pypi":
          results = append(results, validatePyPI(pkg, paths))
        case "npm":
          results = append(results, validateNPM(pkg, paths))
        case "go":
          results = append(results, validateGoModule(pkg, paths))
        case "php":
          results = append(results, validatePHP(pkg, paths))
        case "ruby":
          results = append(results, validateRuby(pkg, paths))
        case "rust":
          results = append(results, validateRust(pkg, paths))
        }
    }
  }

  return results
}
