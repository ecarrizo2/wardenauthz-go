package wardenauth

import (
	"context"
	"fmt"
	"strings"
)

type RolesResource struct{ client *Client }

func (r *RolesResource) Create(ctx context.Context, scopeID string, input CreateRoleInput) (*RoleItem, error) {
	type body struct {
		CreateRoleInput
		ScopeID string `json:"scopeId"`
	}
	var result RoleItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/role", scopeID),
		body{input, scopeID}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RolesResource) BulkCreate(ctx context.Context, scopeID string, inputs []CreateRoleInput) error {
	type body struct {
		CreateRoleInput
		ScopeID string `json:"scopeId"`
	}
	bodies := make([]body, len(inputs))
	for i, inp := range inputs {
		bodies[i] = body{inp, scopeID}
	}
	return r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/role/bulk", scopeID), bodies, nil)
}

func (r *RolesResource) List(ctx context.Context, scopeID string, opts *struct{ Limit, NextToken string }) (*PaginatedResult[RoleItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
	}
	var result PaginatedResult[RoleItem]
	if err := r.client.get(ctx, buildQueryPath(fmt.Sprintf("/v1/scope/%s/role", scopeID), params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RolesResource) GetByID(ctx context.Context, scopeID, roleID string) (*RoleItem, error) {
	var result RoleItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/role/%s", scopeID, roleID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RolesResource) Update(ctx context.Context, scopeID, roleID string, input UpdateRoleInput) (*RoleItem, error) {
	var result RoleItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/role/%s", scopeID, roleID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RolesResource) Delete(ctx context.Context, scopeID, roleID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/role/%s", scopeID, roleID))
}

func (r *RolesResource) BulkDelete(ctx context.Context, scopeID string, ids []string) error {
	return r.client.do(ctx, "DELETE", fmt.Sprintf("/v1/scope/%s/role/bulk", scopeID),
		map[string][]string{"ids": ids}, nil)
}

func (r *RolesResource) Clone(ctx context.Context, targetScopeID, templateRoleID string) (*RoleItem, error) {
	var result RoleItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/role/clone-from/%s", targetScopeID, templateRoleID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func joinQuery(ids []string) string {
	var sb strings.Builder
	for i, id := range ids {
		if i > 0 { sb.WriteString(",") }
		sb.WriteString(id)
	}
	return sb.String()
}
