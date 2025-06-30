package validator

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/Kelcode-Dev/vibe-validator/utils"
)

type pypiMetadata struct {
  Info struct {
    ProjectURL string `json:"project_url"`
    HomePage   string `json:"home_page"`
    PackageURL string `json:"package_url"`
  } `json:"info"`
  Releases map[string][]struct {
    UploadTimeISO string `json:"upload_time_iso_8601"`
  } `json:"releases"`
}

func validatePyPI(packageName string, paths []string) ValidationResult {
  result := ValidationResult{Name: packageName, Source: "pypi", Paths: paths}

  url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)
  resp, err := http.Get(url)
  if err != nil || resp == nil || resp.StatusCode != 200 {
    result.Status = "not_found"
    result.Details = "Not found on PyPI"
    return result
  }
  defer resp.Body.Close()

  var data pypiMetadata
  if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
    result.Status = "investigate"
    result.Details = "Unable to decode PyPI metadata"
    return result
  }

  var oldest time.Time
  for _, versions := range data.Releases {
    for _, release := range versions {
      t, err := time.Parse(time.RFC3339, release.UploadTimeISO)
      if err == nil && (oldest.IsZero() || t.Before(oldest)) {
        oldest = t
      }
    }
  }

  age := time.Since(oldest)
  if age < 30*24*time.Hour {
    result.Status = "investigate"
    result.Details = fmt.Sprintf("Very new package (published %s ago)", utils.HumanDuration(age))
  } else {
    result.Status = "safe"
    result.Details = "-"
  }

  return result
}
