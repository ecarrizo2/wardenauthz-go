package wardenauth

import (
	"context"
)

type AccessResource struct{ client *Client }

func (r *AccessResource) HasAccess(ctx context.Context, input AccessCheckInput) (*AccessCheckResult, error) {
	var result AccessCheckResult
	if err := r.client.post(ctx, "/v1/access/check", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessResource) HasAccessBulk(ctx context.Context, inputs []AccessCheckInput) ([]AccessCheckResult, error) {
	var results []AccessCheckResult
	if err := r.client.post(ctx, "/v1/access/check-bulk", inputs, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *AccessResource) ListPermissions(ctx context.Context, input ListPermissionsInput) (*ListPermissionsResult, error) {
	var result ListPermissionsResult
	if err := r.client.post(ctx, "/v1/access/list-permissions", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessResource) ListRoles(ctx context.Context, input ListRolesInput) (*ListRolesResult, error) {
	var result ListRolesResult
	if err := r.client.post(ctx, "/v1/access/list-roles", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessResource) IssueReceipt(ctx context.Context, input ReceiptIssueInput) (*ReceiptIssueResult, error) {
	var result ReceiptIssueResult
	if err := r.client.post(ctx, "/v1/access/receipt", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessResource) VerifyReceipt(ctx context.Context, input ReceiptVerifyInput) (*ReceiptVerifyResult, error) {
	var result ReceiptVerifyResult
	if err := r.client.post(ctx, "/v1/access/receipt/verify", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessResource) Simulate(ctx context.Context, input SimulateAccessInput) (*SimulateAccessResult, error) {
	var result SimulateAccessResult
	if err := r.client.post(ctx, "/v1/access/simulate", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
