package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// get user email from service account and if that fails then fall back to local email held for gcloud
func GetUserEmail() (email string, err error) {
	// Initialize a context
	ctx := context.Background()

	// Use default credentials
	creds, err := google.FindDefaultCredentials(ctx, oauth2.UserinfoEmailScope)
	if err != nil {
		return "", fmt.Errorf("failed to obtain default credentials: %e", err)
	}

	// Check if credentials are from a service account
	if len(creds.JSON) > 0 {
		var content map[string]interface{}
		if err := json.Unmarshal(creds.JSON, &content); err == nil {
			if email, ok := content["client_email"].(string); ok {
				return email, nil
			}
		}
	}

	// If not, then the credentials are from a user account

	// Create an OAuth2 service using the credentials
	oauth2Service, err := oauth2.NewService(ctx, option.WithTokenSource(creds.TokenSource))
	if err != nil {
		return "", fmt.Errorf("unable to create OAuth2 service: %e", err)
	}

	// Get the authenticated user's email address
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user info: %e", err)
	}

	return userInfo.Email, nil
}
