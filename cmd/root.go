package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const longDescription = `
This CLI allows you to request just-in-time permissions for GCP projects.
Usage example:
  gpc-jit entitlements --project my-project-id
`

var rootCmd = &cobra.Command{
	Use:   "gpc-jit",
	Short: "Request just-in-time GCP permissions",
	Long:  longDescription,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It is called by main.main() and only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
