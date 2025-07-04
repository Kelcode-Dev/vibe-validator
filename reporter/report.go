package reporter

import (
  "fmt"
  "os"
  "sort"
  "strings"
  "text/tabwriter"

  "github.com/Kelcode-Dev/vibe-validator/validator"
)

func PrintReport(results []validator.ValidationResult, verbosity int) {
  if len(results) == 0 {
    fmt.Println("No dependencies found.")
    return
  }

  fmt.Println("\n[vibe-validator] Dependency Vibe Report\n")

  groups := map[string][]validator.ValidationResult{}
  for _, r := range results {
    groups[r.Source] = append(groups[r.Source], r)
  }

  for eco, group := range groups {
    filtered := []validator.ValidationResult{}
    for _, r := range group {
      // Filter by verbosity:
      if verbosity == 0 && r.Status == "safe" {
        continue // default: hide safe
      }
      filtered = append(filtered, r)
    }

    if len(filtered) == 0 {
      continue
    }

    sort.Slice(filtered, func(i, j int) bool {
      return filtered[i].Name < filtered[j].Name
    })

    fmt.Printf("%s:\n", eco)
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "  Status\tName\tDetails\tPath")

    for _, r := range filtered {
      icon := map[string]string{
        "safe":        "[✓]",
        "investigate": "[~]",
        "not_found":   "[✗]",
      }[r.Status]

      detail := r.Details
      if detail == "" {
        detail = "-"
      }
      paths := strings.Join(r.Paths, ", ")
      fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", icon, r.Name, detail, paths)
    }
    w.Flush()
    fmt.Println()
  }
}

