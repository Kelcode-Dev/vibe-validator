package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kelcode-Dev/vibe-validator/utils"
)

type npmMetadata struct {
	Time map[string]string `json:"time"`
}

func validateNPM(packageName string, paths []string) ValidationResult {
	result := ValidationResult{Name: packageName, Source: "npm", Paths: paths}

	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	resp, err := http.Get(url)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		result.Status = "not_found"
		result.Details = "Not found on npm"
		return result
	}
	defer resp.Body.Close()

	var data npmMetadata
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		result.Status = "investigate"
		result.Details = "Unable to decode npm metadata"
		return result
	}

	createdAt := data.Time["created"]
	t, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		result.Status = "investigate"
		result.Details = "Invalid publish timestamp"
		return result
	}

	age := time.Since(t)
	if age < 30*24*time.Hour {
		result.Status = "investigate"
		result.Details = fmt.Sprintf("Very new package (published %s ago)", utils.HumanDuration(age))
	} else {
		result.Status = "safe"
		result.Details = "-"
	}

	return result
}
