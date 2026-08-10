package wardenauth

import (
	"context"
	"fmt"
)

type AccessPoliciesResource struct{ client *Client }

func (r *AccessPoliciesResource) Create(ctx context.Context, scopeID string, input CreateAccessPolicyInput) (*AccessPolicyItem, error) {
	type body struct {
		CreateAccessPolicyInput
		ScopeID string `json:"scopeId"`
	}
	var result AccessPolicyItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/access-policy", scopeID),
		body{input, scopeID}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessPoliciesResource) ListByScope(ctx context.Context, scopeID string) ([]AccessPolicyItem, error) {
	type output struct{ Items []AccessPolicyItem `json:"items"` }
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/access-policy", scopeID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *AccessPoliciesResource) ListBySubject(ctx context.Context, scopeID, subjectID string) ([]AccessPolicyItem, error) {
	type output struct{ Items []AccessPolicyItem `json:"items"` }
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/access-policy/subject/%s", scopeID, subjectID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *AccessPoliciesResource) GetByID(ctx context.Context, scopeID, policyID string) (*AccessPolicyItem, error) {
	var result AccessPolicyItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/access-policy/%s", scopeID, policyID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessPoliciesResource) Update(ctx context.Context, scopeID, policyID string, input UpdateAccessPolicyInput) (*AccessPolicyItem, error) {
	var result AccessPolicyItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/access-policy/%s", scopeID, policyID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccessPoliciesResource) Delete(ctx context.Context, scopeID, policyID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/access-policy/%s", scopeID, policyID))
}

func (r *AccessPoliciesResource) ImportCSV(ctx context.Context, scopeID, csvContent string) (*ImportResult, error) {
	var result ImportResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/access-policy/import-csv", scopeID),
		map[string]string{"csv": csvContent}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
