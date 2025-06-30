package validator

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/Kelcode-Dev/vibe-validator/utils"
)

type goModuleInfo struct {
  Version string    `json:"Version"`
  Time    time.Time `json:"Time"`
}

func validateGoModule(module string, paths []string) ValidationResult {
  result := ValidationResult{Name: module, Source: "go", Paths: paths}

  url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", module)
  resp, err := http.Get(url)
  if err != nil || resp == nil || resp.StatusCode != 200 {
    result.Status = "not_found"
    result.Details = "Not found in Go proxy"
    return result
  }
  defer resp.Body.Close()

  var info goModuleInfo
  if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
    result.Status = "unknown"
    result.Details = "Unable to parse module metadata"
    return result
  }

  age := time.Since(info.Time)
  if age < 30*24*time.Hour {
    result.Status = "investigate"
    result.Details = fmt.Sprintf("Recently added (%s)", utils.HumanDuration(age))
  } else {
    result.Status = "safe"
    result.Details = "-"
  }

  return result
}

