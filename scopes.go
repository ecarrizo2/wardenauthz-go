package wardenauth

import (
	"context"
	"fmt"
)

type ScopesResource struct{ client *Client }

func (r *ScopesResource) Create(ctx context.Context, input CreateScopeInput) (*ScopeItem, error) {
	var result ScopeItem
	if err := r.client.post(ctx, "/v1/scope", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ScopesResource) List(ctx context.Context, opts *struct {
	Limit     string
	NextToken string
	Type      ScopeType
}) (*PaginatedResult[ScopeItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
		if opts.Type != "" { params["type"] = string(opts.Type) }
	}
	var result PaginatedResult[ScopeItem]
	if err := r.client.get(ctx, buildQueryPath("/v1/scope", params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ScopesResource) GetByID(ctx context.Context, scopeID string) (*ScopeItem, error) {
	var result ScopeItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s", scopeID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ScopesResource) Update(ctx context.Context, scopeID string, input UpdateScopeInput) (*ScopeItem, error) {
	var result ScopeItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ScopesResource) Delete(ctx context.Context, scopeID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s", scopeID))
}

func (r *ScopesResource) ApplyManifest(ctx context.Context, scopeID string, input ApplyScopeManifestInput) (*ApplyScopeManifestResult, error) {
	var result ApplyScopeManifestResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/apply", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ScopesResource) Export(ctx context.Context, scopeID string) (map[string]any, error) {
	var result map[string]any
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/export", scopeID), &result); err != nil {
		return nil, err
	}
	return result, nil
}
