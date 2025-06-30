package cmd

import (
  "fmt"
  "os"
  "time"

  "github.com/spf13/cobra"
  "github.com/Kelcode-Dev/vibe-validator/scanner"
  "github.com/Kelcode-Dev/vibe-validator/validator"
  "github.com/Kelcode-Dev/vibe-validator/reporter"
  "github.com/briandowns/spinner"
)

var includeLockfiles bool
var includeVendor bool

func init() {
  rootCmd.Flags().BoolVar(&includeLockfiles, "include-lockfiles", false, "Scan npm/yarn lockfiles for all dependencies")
  rootCmd.Flags().BoolVar(&includeVendor, "include-vendor", false, "Include vendor directories in scanning (slow!)")
}

var rootCmd = &cobra.Command{
  Use:   "vibe-validator [path]",
  Short: "Scan project dependencies for sketchy vibes",
  Args:  cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    path := args[0]
    fmt.Println("[oo] Scanning:", path)

    s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
    s.Suffix = " scanning dependencies..."
    s.Start()
    defer s.Stop()

    pyDeps, npmDeps, goDeps := scanner.ScanDependencies(path, includeLockfiles, includeVendor)
    results := validator.ValidatePackages(pyDeps, npmDeps, goDeps)

    s.Stop() // stop spinner before printing
    reporter.PrintReport(results)
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println("❌ Error:", err)
    os.Exit(1)
  }
}
