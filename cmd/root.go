package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/Kelcode-Dev/vibe-validator/scanner"
  "github.com/Kelcode-Dev/vibe-validator/validator"
  "github.com/Kelcode-Dev/vibe-validator/reporter"
)

var includeLockfiles bool

func init() {
  rootCmd.Flags().BoolVar(&includeLockfiles, "include-lockfiles", false, "Scan npm/yarn lockfiles for all dependencies")
}

var rootCmd = &cobra.Command{
  Use:   "vibe-validator [path]",
  Short: "Scan project dependencies for sketchy vibes",
  Args:  cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    path := args[0]
    fmt.Println("[oo] Scanning:", path)

    pyDeps, npmDeps, goDeps := scanner.ScanDependencies(path, includeLockfiles)
    results := validator.ValidatePackages(pyDeps, npmDeps, goDeps)
    reporter.PrintReport(results)
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println("❌ Error:", err)
    os.Exit(1)
  }
}
