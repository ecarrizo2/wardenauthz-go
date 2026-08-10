package wardenauth

import (
	"context"
	"fmt"
)

type IntegrationsResource struct{ client *Client }

func (r *IntegrationsResource) Create(ctx context.Context, scopeID string, input CreateIntegrationInput) (*Integration, error) {
	var result Integration
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/integration", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *IntegrationsResource) List(ctx context.Context, scopeID string) ([]Integration, error) {
	type output struct {
		Items []Integration `json:"items"`
	}
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/integration", scopeID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *IntegrationsResource) GetByID(ctx context.Context, scopeID, id string) (*Integration, error) {
	var result Integration
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/integration/%s", scopeID, id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *IntegrationsResource) Update(ctx context.Context, scopeID, id string, input UpdateIntegrationInput) (*Integration, error) {
	var result Integration
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/integration/%s", scopeID, id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *IntegrationsResource) Delete(ctx context.Context, scopeID, id string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/integration/%s", scopeID, id))
}
