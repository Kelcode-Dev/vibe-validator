package validator

// ValidationResult holds dependency check results with multiple paths
type ValidationResult struct {
  Name    string
  Source  string   // "npm", "pypi", "go"
  Status  string   // "safe", "investigate", "not_found"
  Details string
  Paths   []string `json:"paths"`
}

// ValidatePackages runs validation on py, npm, and go dependencies with multiple paths
func ValidatePackages(
  pyDeps map[string][]string,
  npmDeps map[string][]string,
  goDeps map[string][]string,
) []ValidationResult {
  var results []ValidationResult

  for pkg, paths := range pyDeps {
    results = append(results, validatePyPI(pkg, paths))
  }
  for pkg, paths := range npmDeps {
    results = append(results, validateNPM(pkg, paths))
  }
  for pkg, paths := range goDeps {
    results = append(results, validateGoModule(pkg, paths))
  }

  return results
}
