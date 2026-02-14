package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/felixgborrego/gpc-pam-jit/pkg/config"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"
	"github.com/felixgborrego/gpc-pam-jit/pkg/slack"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request an entitlement",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID, err := cmd.Flags().GetString("project")
		if err != nil {
			return
		}

		folderID, err := cmd.Flags().GetString("folder")
		if err != nil {
			return
		}

		if (projectID == "" && folderID == "") || (projectID != "" && folderID != "") {
			return
		}

		resourceID := projectID
		resourceType := pamjit.ResourceTypeProject
		if resourceID == "" {
			resourceID = folderID
			resourceType = pamjit.ResourceTypeFolder
		}

		options := &pamjit.RequestOptions{
			EntitlementID: args[0],
			Location:      cmd.Flag("location").Value.String(),
			ResourceID:    resourceID,
			Justification: cmd.Flag("justification").Value.String(),
			Duration:      cmd.Flag("duration").Value.String(),
		}

		pam, err := pamjit.NewPamJitClient(context.Background(), options.ResourceID, options.Location, resourceType)
		if err != nil {
			log.Fatalf("unable to use GCP JIT service: %v", err)
		}
		link, err := pam.RequestGrant(cmd.Context(), options.EntitlementID, options.Justification, options.Duration)
		if err != nil {
			fmt.Printf("Error requesting entitlement: %v\n", err)
		} else {
			if link != "" {

				cfg, _ := config.LoadConfig()

				// only attempt to send to Slack if config is set
				if cfg.Slack.Token != "" && cfg.Slack.Channel != "" {
					// send the link to Slack and if it fails then display the link
					err = slack.SendSlackMessage(cfg, options, link)
					if err != nil {
						fmt.Printf("Link to request: %s\n", link)
					}
				} else {
					fmt.Printf("Link to request: %s\n", link)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("project", "p", "", "Project ID")
	requestCmd.Flags().StringP("folder", "f", "", "Folder ID")
	requestCmd.Flags().StringP("location", "l", "global", "Location")
	requestCmd.Flags().StringP("justification", "j", "", "Justification")
	requestCmd.Flags().StringP("duration", "d", "", "Duration (defaults to maximum)")

	requestCmd.MarkFlagsOneRequired("project", "folder")
	requestCmd.MarkFlagsMutuallyExclusive("project", "folder")
}
