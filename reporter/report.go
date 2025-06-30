package reporter

import (
  "fmt"
  "os"
  "sort"
  "strings"
  "text/tabwriter"

  "github.com/Kelcode-Dev/vibe-validator/validator"
)

func PrintReport(results []validator.ValidationResult) {
  if len(results) == 0 {
    fmt.Println("No dependencies found.")
    return
  }

  fmt.Println("\n[vibe-validator] Dependency Vibe Report\n")

  groups := map[string][]validator.ValidationResult{}
  for _, r := range results {
    groups[r.Source] = append(groups[r.Source], r)
  }

  ecosystems := []string{"pypi", "npm", "go"}
  for _, eco := range ecosystems {
    group, ok := groups[eco]
    if !ok {
      continue
    }

    sort.Slice(group, func(i, j int) bool {
      return group[i].Name < group[j].Name
    })

    fmt.Printf("%s:\n", eco)
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "  Status\tName\tDetails\tPath")
    for _, r := range group {
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
