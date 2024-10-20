package cmd

import (
	"context"
	"log"

	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"

	"github.com/spf13/cobra"
)

var listEntitlementCmd = &cobra.Command{
	Use:   "entitlements",
	Short: "List entitlements",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		location, _ := cmd.Flags().GetString("location")

		pam, err := pamjit.NewPamJitClient(context.Background(), projectID, location)
		if err != nil {
			log.Fatalf("unable to use GCP JIT service: %v", err)
		}
		err = pam.ShowEntitlements(context.Background())
		if err != nil {
			log.Fatalf("Error listing entitlements: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listEntitlementCmd)

	listEntitlementCmd.Flags().StringP("project", "p", "", "Project ID")
	listEntitlementCmd.Flags().StringP("location", "l", "global", "Location")
	_ = listEntitlementCmd.MarkFlagRequired("project")
}
