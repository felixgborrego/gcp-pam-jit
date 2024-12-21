package pamjit

import (
	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"regexp"
	"strconv"
	"strings"
)

type RequestOptions struct {
	EntitlementID string
	ProjectID     string
	Location      string
	Justification string
	Duration      string
}

func (c *Client) RequestGrant(ctx context.Context, entitlementId, justification string, duration string) (string, error) {
	entitlementFqn := fmt.Sprintf("%s/entitlements/%s", c.parent(), entitlementId)

	reqEntitlement := &privilegedaccessmanagerpb.GetEntitlementRequest{
		Name: entitlementFqn,
	}

	entitlement, err := c.gcpClient.GetEntitlement(ctx, reqEntitlement)
	if err != nil {
		return "", fmt.Errorf("entitlement not found: %v", err)
	}

	// convert duration string to google.protobuf.Duration
	var requestedDuration *durationpb.Duration
	if duration != "" {
		var err error
		requestedDuration, err = parseDurationProto(duration)
		if err != nil {
			return "", fmt.Errorf("invalid duration: %w", err)
		}
	} else {
		requestedDuration = entitlement.MaxRequestDuration
	}

	fmt.Printf("Requesting entitlement %s in project %s for %s\n", entitlementId, c.projectID, requestedDuration)

	req := &privilegedaccessmanagerpb.CreateGrantRequest{
		Parent: entitlementFqn,
		Grant: &privilegedaccessmanagerpb.Grant{
			Name:              entitlementId,
			RequestedDuration: requestedDuration,
			Justification: &privilegedaccessmanagerpb.Justification{
				Justification: &privilegedaccessmanagerpb.Justification_UnstructuredJustification{
					UnstructuredJustification: justification,
				},
			},
		},
	}

	grant, err := c.gcpClient.CreateGrant(ctx, req)

	if err != nil {
		return "", err
	}

	fmt.Printf("Grant request sent: %s\n", grant.GetState().String())

	link := fmt.Sprintf("https://console.cloud.google.com/iam-admin/pam/grants/approvals?project=%s", c.projectID)

	// only return link if the request requires approval
	if grant.GetState().String() == "APPROVAL_AWAITED" {
		return link, nil
	} else {
		return "", nil
	}
}

// parseDurationProto converts a duration string like "30s", "5m", or "2h" to *duration.Duration
func parseDurationProto(durationStr string) (*durationpb.Duration, error) {
	durationStr = strings.TrimSpace(durationStr)
	if durationStr == "" {
		return nil, fmt.Errorf("duration string is empty")
	}

	// use a regular expression to extract the numeric value and unit
	re := regexp.MustCompile(`^(\d+)([smh])$`)
	matches := re.FindStringSubmatch(durationStr)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid duration string format: %s", durationStr)
	}

	// parse the numeric value
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid duration value: %s", matches[1])
	}

	// calculate the duration in seconds based on the unit
	var seconds int64
	switch matches[2] {
	case "s":
		seconds = int64(value)
	case "m":
		seconds = int64(value) * 60
	case "h":
		seconds = int64(value) * 60 * 60
	default:
		return nil, fmt.Errorf("invalid duration unit: %s", matches[2])
	}

	return &durationpb.Duration{Seconds: seconds}, nil
}
