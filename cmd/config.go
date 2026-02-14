package cmd

import (
	"fmt"

	"github.com/felixgborrego/gpc-pam-jit/pkg/config"
	"github.com/spf13/cobra"
)

// configCmd defines the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the GCP JIT service",
}

// slackCmd represents the slack subcommand to configure the Slack integration
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Configure the Slack integration",
	Run: func(cmd *cobra.Command, args []string) {
		slackConfig(cmd)
	},
}

func slackConfig(cmd *cobra.Command) {
	channel, _ := cmd.Flags().GetString("channel")
	token, _ := cmd.Flags().GetString("token")

	// Save the configuration
	conf := config.Config{
		Slack: config.SlackConfig{
			Channel: channel,
			Token:   token,
		},
	}

	if err := config.SaveConfig(&conf); err != nil {
		fmt.Println("Error saving the configuration:", err)
		return
	}
	fmt.Println("Configuration saved")
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(slackCmd)

	slackCmd.Flags().StringP("channel", "c", "", "Slack channel to send messages to")
	slackCmd.Flags().StringP("token", "t", "", "Slack token")

	slackCmd.MarkFlagRequired("channel")
	slackCmd.MarkFlagRequired("token")
}
