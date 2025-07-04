package validator

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"
)

type cratesResponse struct {
  Crate struct {
    CreatedAt string `json:"created_at"` // ISO8601
  } `json:"crate"`
}

func validateRust(crateName string, paths []string) ValidationResult {
  result := ValidationResult{
    Name:   crateName,
    Source: "rust",
    Paths:  paths,
  }

  url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", crateName)
  resp, err := http.Get(url)
  if err != nil || resp == nil || resp.StatusCode != 200 {
    result.Status = "not_found"
    result.Details = "Not found on crates.io"
    return result
  }
  defer resp.Body.Close()

  var data cratesResponse
  if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
    result.Status = "investigate"
    result.Details = "Unable to decode crates.io metadata"
    return result
  }

  t, err := time.Parse(time.RFC3339, data.Crate.CreatedAt)
  if err != nil {
    result.Status = "investigate"
    result.Details = "Invalid publish timestamp"
    return result
  }

  age := time.Since(t)
  if age < 30*24*time.Hour {
    result.Status = "investigate"
    result.Details = fmt.Sprintf("Very new package (published %s ago)", age.Round(time.Hour*24))
  } else {
    result.Status = "safe"
    result.Details = "-"
  }

  return result
}
