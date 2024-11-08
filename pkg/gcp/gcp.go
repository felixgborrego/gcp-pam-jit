package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/oauth2/google"
)

// get user email from service account and if that fails then fall back to local email held for gcloud
func GetUserEmail() (email string, err error) {
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
