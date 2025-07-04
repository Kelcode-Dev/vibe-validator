package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type packagistResponse struct {
	Packages map[string][]struct {
		Time string `json:"time"` // ISO8601 timestamp of the version release
	} `json:"packages"`
}

func validatePHP(packageName string, paths []string) ValidationResult {
	result := ValidationResult{
		Name:   packageName,
		Source: "php",
		Paths:  paths,
	}

	url := fmt.Sprintf("https://repo.packagist.org/p/%s.json", packageName)
	resp, err := http.Get(url)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		result.Status = "not_found"
		result.Details = "Not found on Packagist"
		return result
	}
	defer resp.Body.Close()

	var data packagistResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		result.Status = "investigate"
		result.Details = "Unable to decode Packagist metadata"
		return result
	}

	versions, found := data.Packages[packageName]
	if !found || len(versions) == 0 {
		result.Status = "not_found"
		result.Details = "No versions found on Packagist"
		return result
	}

	// Find oldest release time
	var oldest time.Time
	for _, v := range versions {
		t, err := time.Parse(time.RFC3339, v.Time)
		if err == nil && (oldest.IsZero() || t.Before(oldest)) {
			oldest = t
		}
	}

	age := time.Since(oldest)
	if age < 30*24*time.Hour {
		result.Status = "investigate"
		result.Details = fmt.Sprintf("Very new package (published %s ago)", age.Round(time.Hour*24))
	} else {
		result.Status = "safe"
		result.Details = "-"
	}

	return result
}
