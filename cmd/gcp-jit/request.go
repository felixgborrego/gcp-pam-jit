package main

import (
	"context"
	"fmt"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"
	"log"

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

		pam, err := pamjit.NewPamJitClient(context.Background(), projectID, location)
		if err != nil {
			log.Fatalf("unable to use GCP JIT service: %v", err)
		}
		err = pam.RequestGrant(cmd.Context(), entitlementID, justification)
		if err != nil {
			fmt.Printf("Error requesting entitlement: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().String("project", "", "Project ID")
	requestCmd.Flags().String("location", "global", "Location")
	requestCmd.Flags().String("justification", "", "Justification")

	requestCmd.MarkFlagRequired("project")
	requestCmd.MarkFlagRequired("justification")
}
