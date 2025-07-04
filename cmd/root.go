package cmd

import (
  "fmt"
  "time"
  "os"

  "github.com/spf13/cobra"
  "github.com/Kelcode-Dev/vibe-validator/scanner"
  "github.com/Kelcode-Dev/vibe-validator/validator"
  "github.com/Kelcode-Dev/vibe-validator/reporter"
  "github.com/briandowns/spinner"
)

var (
  includeLockfiles bool
  includeVendor bool
  verbosity int = 0
)

func init() {
  rootCmd.Flags().BoolVar(&includeLockfiles, "include-lockfiles", false, "Scan npm/yarn lockfiles for all dependencies")
  rootCmd.Flags().BoolVar(&includeVendor, "include-vendor", false, "Include vendor directories in scanning (slow!)")
  rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity level")
}

var rootCmd = &cobra.Command{
  Use:   "vibe-validator [path]",
  Short: "Scan project dependencies for sketchy vibes",
  Args:  cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    path := args[0]
    //create a cli logo for the top of the output
    fmt.Println("[oo] Scanning:", path)

    opts := scanner.ScanOptions{
      IncludeLockfiles: includeLockfiles,
      IncludeVendor:    includeVendor,
      Verbosity:        verbosity,
    }

    s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
    if verbosity < 2 {
      s.Suffix = " scanning dependencies..."
      s.Start()
      defer s.Stop()
    }

    deps, err := scanner.ScanDependencies(path, opts)
    if verbosity < 2 {
      s.Stop()
    }
    if err != nil {
      fmt.Printf("❌ Scan failed: %v\n", err)
      os.Exit(1)
    }

    v := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
    v.Suffix = " Validating dependencies..."
    v.Start()
    defer v.Stop()

    results := validator.ValidatePackages(deps)
    v.Stop()

    fmt.Println("Validation complete, prepping report...")
    reporter.PrintReport(results, verbosity)
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println("❌ Error:", err)
    os.Exit(1)
  }
}
