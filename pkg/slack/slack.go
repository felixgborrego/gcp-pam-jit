package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/oauth2/google"
	"github.com/felixgborrego/gpc-pam-jit/pkg/config"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"
	"github.com/slack-go/slack"
)

func SendSlackMessage(cfg *config.Config, options *pamjit.RequestOptions, link string) (error) {
	api := slack.New(cfg.Slack.Token)
	email, _ := getUserEmail()

	// Use Slack Block Kit for better formatting
	blocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", ":lock: PAM Request", false, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Entitlement:*\n%s", options.EntitlementID), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Resource:*\n%s", options.ProjectID), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Requested by:*\n%s", email), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Duration:*\n%s", options.Duration), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Justification:*\n%s", options.Justification), false, false),
			nil,
			nil,
		),
		slack.NewActionBlock(
			"",
			slack.NewButtonBlockElement("", "Approve or Deny", slack.NewTextBlockObject("plain_text", "Approve or Deny", false, false)).WithURL(link),
		),
	}

	// Send the message with blocks
	_, _, err := api.PostMessage(cfg.Slack.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("error sending message to Slack: %e", err)
	}

	fmt.Println("Sent request to Slack")

	return nil
}

// get user email from service account and if that fails then fall back to local email held for gcloud
func getUserEmail() (email string, err error) {
	ctx := context.Background()
	credential, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting default credentials: %e", err)
	}

	content := map[string]interface{}{}

	json.Unmarshal(credential.JSON, &content)

	if content["client_email"] != nil {
		return content["client_email"].(string), nil
	} else {
		email, err := getGcloudUserEmail()
		if err != nil {
			return "", fmt.Errorf("error getting user email: %e", err)
		}
		return email, nil
	}
}

func getGcloudUserEmail() (string, error) {
	cmd := exec.Command("gcloud", "auth", "list", "--format=value(account)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing gcloud command: %w", err)
	}

	email := strings.TrimSpace(string(output))
	return email, nil
}
