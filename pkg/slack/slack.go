package slack

import (
	"fmt"

	"github.com/felixgborrego/gpc-pam-jit/pkg/config"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"
	"github.com/felixgborrego/gpc-pam-jit/pkg/gcp"
	"github.com/slack-go/slack"
)

func SendSlackMessage(cfg *config.Config, options *pamjit.RequestOptions, link string) error {
	api := slack.New(cfg.Slack.Token)
	email, _ := gcp.GetUserEmail()

	// Use Slack Block Kit for better formatting
	blocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", ":lock: PAM Request", false, false),
		),
		// wait for voting on format...
		// slack.NewSectionBlock(
		// 	slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Entitlement:* `%s`\n*Resource:* `%s`\n*Requested by:* `%s`\n*Duration:* `%s`\n*Justification:* `%s`", options.EntitlementID, options.ProjectID, email, options.Duration, options.Justification), false, false),
		// 	nil,
		// 	nil,
		// ),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s* has requested *%s* on resource *%s* for *%s*, with the justification\n>%s", email, options.EntitlementID, options.ProjectID, options.Duration, options.Justification), false, false),
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
	_, _, err := api.PostMessage(cfg.Slack.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("error sending message to Slack: %e", err)
	}

	fmt.Println("Sent request to Slack")

	return nil
}
