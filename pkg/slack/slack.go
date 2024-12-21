package slack

import (
	"fmt"

	"github.com/felixgborrego/gpc-pam-jit/pkg/config"
	"github.com/felixgborrego/gpc-pam-jit/pkg/gcp"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"
	"github.com/slack-go/slack"
)

func SendSlackMessage(cfg *config.Config, options *pamjit.RequestOptions, link string) error {
	api := slack.New(cfg.Slack.Token)
	email, err := gcp.GetUserEmail()
	if email == "" {
		fmt.Printf("Unable to retrieve your email but will send the request to Slack anyway\n%s\n", err)
	}

	// Use Slack Block Kit for better formatting
	blocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", ":lock: PAM Request", false, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Entitlement: *%s*\nResource: *%s*\nRequested by: *%s*\nDuration: *%s*\nJustification: *%s*", options.EntitlementID, options.ProjectID, email, options.Duration, options.Justification), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Please review and approve <%s|here>", link), false, false),
			nil,
			nil,
		),
	}

	// send the message with blocks
	_, _, err = api.PostMessage(cfg.Slack.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("error sending message to Slack: %e", err)
	}

	fmt.Println("Sent request to Slack")

	return nil
}
