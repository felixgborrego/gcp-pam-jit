package pamjit

import (
	"context"
	"fmt"
	"strings"
	"time"

	privilegedaccessmanagerpb "cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"google.golang.org/api/iterator"
)

const (
	pageSize = 10
)

// ShowEntitlements lists the entitlements for a given project ID and location.
func (c *Client) ShowEntitlements(ctx context.Context) error {
	fmt.Printf("Your current GCP user has the following entitlements for the project %s and location %s:\n\n", c.projectID, c.location)

	req := &privilegedaccessmanagerpb.ListEntitlementsRequest{
		Parent: c.parent(),
	}

	var allEntitlements []*privilegedaccessmanagerpb.Entitlement
	it := c.gcpClient.ListEntitlements(ctx, req)
	pager := iterator.NewPager(it, pageSize, "")
	for {
		var entitlements []*privilegedaccessmanagerpb.Entitlement
		nextPageToken, err := pager.NextPage(&entitlements)
		if err != nil {
			return fmt.Errorf("failed to retrieve entitlements: %w", err)
		}
		allEntitlements = append(allEntitlements, entitlements...)

		if nextPageToken == "" {
			break
		}
	}

	if len(allEntitlements) == 0 {
		fmt.Println("No entitlements found.")
		return nil
	}

	for _, entitlement := range allEntitlements {
		printEntitlement(entitlement)
	}

	return nil
}

// printEntitlement prints the entitlement details in a human-friendly format.
func printEntitlement(entitlement *privilegedaccessmanagerpb.Entitlement) {
	maxRequestDuration := time.Duration(entitlement.MaxRequestDuration.Seconds) * time.Second
	entitlementName := entitlementNameFromFullName(entitlement.Name)
	fmt.Printf("🛡️ Entitlement: %s (%s)\n", entitlementName, maxRequestDuration)

	roles := extractRoles(entitlement)
	fmt.Printf("    Granted Roles: %s\n", strings.Join(roles, ", "))

	printApprovalWorkflow(4, entitlement)
}

// entitlementNameFromFullName extracts the entitlement name from its full resource name.
func entitlementNameFromFullName(fullName string) string {
	parts := strings.Split(fullName, "/")
	return parts[len(parts)-1]
}

// extractRoles extracts the roles from the entitlement's GCP IAM access bindings.
func extractRoles(entitlement *privilegedaccessmanagerpb.Entitlement) []string {
	var roles []string
	for _, binding := range entitlement.PrivilegedAccess.GetGcpIamAccess().RoleBindings {
		roles = append(roles, binding.Role)
	}
	return roles
}

// printApprovalWorkflow prints the approval workflow details in a human-friendly format.
func printApprovalWorkflow(leftPadding int, entitlement *privilegedaccessmanagerpb.Entitlement) {
	if entitlement.ApprovalWorkflow == nil {
		PrintLine(leftPadding, "No manual approval required\n")
		return
	}

	approvalWorkflow := entitlement.ApprovalWorkflow.GetApprovalWorkflow()
	if approvalWorkflow == nil {
		return
	}

	manualApprovalWorkflow, ok := approvalWorkflow.(*privilegedaccessmanagerpb.ApprovalWorkflow_ManualApprovals)
	if !ok {
		PrintLine(leftPadding, "No manual approval required\n")
		return
	}

	approvers := extractApprovers(manualApprovalWorkflow.ManualApprovals)
	PrintLine(leftPadding, "Approval required by: %s\n", strings.Join(approvers, ", "))
}

// extractApprovers extracts the approvers' principals from the manual approval workflow.
func extractApprovers(manualApprovals *privilegedaccessmanagerpb.ManualApprovals) []string {
	var approvePrincipals []string
	for _, step := range manualApprovals.Steps {
		for _, approver := range step.Approvers {
			approvePrincipals = append(approvePrincipals, strings.Join(approver.Principals, ", "))
		}
	}
	return approvePrincipals
}
