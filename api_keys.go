package wardenauth

import (
	"context"
	"fmt"
)

type APIKeysResource struct{ client *Client }

func (r *APIKeysResource) Create(ctx context.Context, scopeID, keyType string, input CreateApiKeyInput) (*ApiKeyCreatedItem, error) {
	var result ApiKeyCreatedItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s", scopeID, keyType), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *APIKeysResource) List(ctx context.Context, scopeID, keyType string) ([]ApiKeyItem, error) {
	type output struct{ Items []ApiKeyItem `json:"items"` }
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s", scopeID, keyType), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *APIKeysResource) GetByID(ctx context.Context, scopeID, keyType, keyID string) (*ApiKeyItem, error) {
	var result ApiKeyItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s/%s", scopeID, keyType, keyID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *APIKeysResource) Delete(ctx context.Context, scopeID, keyType, keyID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s/%s", scopeID, keyType, keyID))
}

func (r *APIKeysResource) Rotate(ctx context.Context, scopeID, keyType, keyID string) (*ApiKeyRotationResult, error) {
	var result ApiKeyRotationResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s/%s/rotate", scopeID, keyType, keyID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *APIKeysResource) RevealRotation(ctx context.Context, scopeID, keyType, keyID, rotationRef string) (*ApiKeyRevealRotationResult, error) {
	var result ApiKeyRevealRotationResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s/%s/rotation/reveal", scopeID, keyType, keyID),
		map[string]string{"rotationRef": rotationRef}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *APIKeysResource) UpdateAutoRotation(ctx context.Context, scopeID, keyType, keyID string, input UpdateApiKeyAutoRotationInput) (*ApiKeyItem, error) {
	var result ApiKeyItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/api-key/%s/%s/auto-rotation", scopeID, keyType, keyID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
