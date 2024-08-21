package pamjit

import (
	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"context"
	"fmt"
)

func (c *Client) RequestGrant(ctx context.Context, entitlementId, justification string) error {
	entitlementFqn := fmt.Sprintf("%s/entitlements/%s", c.parent(), entitlementId)

	reqEntitlement := &privilegedaccessmanagerpb.GetEntitlementRequest{
		Name: entitlementFqn,
	}

	entitlement, err := c.gcpClient.GetEntitlement(ctx, reqEntitlement)
	if err != nil {
		return fmt.Errorf("entitlement not found: %v", err)
	}
	fmt.Printf("Requesting entitlement %s in project %s for %s\n", entitlementId, c.projectID, entitlement.MaxRequestDuration)

	req := &privilegedaccessmanagerpb.CreateGrantRequest{
		Parent: entitlementFqn,
		Grant: &privilegedaccessmanagerpb.Grant{
			Name:              entitlementId,
			RequestedDuration: entitlement.MaxRequestDuration,
			Justification: &privilegedaccessmanagerpb.Justification{
				Justification: &privilegedaccessmanagerpb.Justification_UnstructuredJustification{
					UnstructuredJustification: justification,
				},
			},
		},
	}

	grant, err := c.gcpClient.CreateGrant(ctx, req)

	if err != nil {
		return err
	}

	fmt.Printf("Granted request sent! %s\n", grant.GetState().String())

	return nil
}
