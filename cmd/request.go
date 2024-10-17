package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request an entitlement",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entitlementID := args[0]
		projectID, _ := cmd.Flags().GetString("project")
		location, _ := cmd.Flags().GetString("location")
		justification, _ := cmd.Flags().GetString("justification")
		duration, _ := cmd.Flags().GetString("duration")

		pam, err := pamjit.NewPamJitClient(context.Background(), projectID, location)
		if err != nil {
			log.Fatalf("unable to use GCP JIT service: %v", err)
		}
		err = pam.RequestGrant(cmd.Context(), entitlementID, justification, duration)
		if err != nil {
			fmt.Printf("Error requesting entitlement: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("project", "p", "", "Project ID")
	requestCmd.Flags().StringP("location", "l", "global", "Location")
	requestCmd.Flags().StringP("justification", "j", "", "Justification")
	requestCmd.Flags().StringP("duration", "d", "", "Duration (defaults to maximum)")

	requestCmd.MarkFlagRequired("project")
	requestCmd.MarkFlagRequired("justification")
}
