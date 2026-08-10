package wardenauth

import (
	"context"
	"fmt"
)

type PermissionsResource struct{ client *Client }

func (r *PermissionsResource) Create(ctx context.Context, scopeID string, input CreatePermissionInput) (*PermissionItem, error) {
	type body struct {
		CreatePermissionInput
		ScopeID string `json:"scopeId"`
	}
	var result PermissionItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/permission", scopeID),
		body{input, scopeID}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PermissionsResource) BulkCreate(ctx context.Context, scopeID string, inputs []CreatePermissionInput) error {
	type body struct {
		CreatePermissionInput
		ScopeID string `json:"scopeId"`
	}
	bodies := make([]body, len(inputs))
	for i, inp := range inputs {
		bodies[i] = body{inp, scopeID}
	}
	return r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/permission/bulk", scopeID), bodies, nil)
}

func (r *PermissionsResource) List(ctx context.Context, scopeID string, opts *struct{ Limit, NextToken string }) (*PaginatedResult[PermissionItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
	}
	var result PaginatedResult[PermissionItem]
	if err := r.client.get(ctx, buildQueryPath(fmt.Sprintf("/v1/scope/%s/permission", scopeID), params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PermissionsResource) GetByID(ctx context.Context, scopeID, permissionID string) (*PermissionItem, error) {
	var result PermissionItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/permission/%s", scopeID, permissionID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PermissionsResource) Update(ctx context.Context, scopeID, permissionID string, input UpdatePermissionInput) (*PermissionItem, error) {
	var result PermissionItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/permission/%s", scopeID, permissionID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PermissionsResource) Delete(ctx context.Context, scopeID, permissionID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/permission/%s", scopeID, permissionID))
}

func (r *PermissionsResource) BulkDelete(ctx context.Context, scopeID string, ids []string) error {
	return r.client.do(ctx, "DELETE", fmt.Sprintf("/v1/scope/%s/permission/bulk", scopeID),
		map[string][]string{"ids": ids}, nil)
}

func (r *PermissionsResource) ImportCSV(ctx context.Context, scopeID, csvContent string) (*ImportResult, error) {
	var result ImportResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/permission/import-csv", scopeID),
		map[string]string{"csv": csvContent}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

