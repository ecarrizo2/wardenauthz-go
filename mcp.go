package wardenauth

import (
	"context"
	"fmt"
	"net/url"
)

// TrustTier is the trust level a user may adopt when connecting to an MCP server.
type TrustTier string

const (
	TrustTierLow    TrustTier = "low"
	TrustTierMedium TrustTier = "medium"
	TrustTierHigh   TrustTier = "high"
)

// McpUserAssignment is an admin-administered grant giving a user access to one MCP server,
// capped at a maximum trust tier they may self-adopt at consent time.
type McpUserAssignment struct {
	SubjectID           string              `json:"subjectId"`
	ServerKey           string              `json:"serverKey"`
	MaxTier             TrustTier           `json:"maxTier"`
	ResourceConstraints map[string][]string `json:"resourceConstraints,omitempty"`
	AssignedBy          string              `json:"assignedBy"`
	CreatedAt           float64             `json:"createdAt"`
	UpdatedAt           float64             `json:"updatedAt"`
}

// MCPResource administers per-user MCP server access assignments (the "Humans" access matrix).
type MCPResource struct{ client *Client }

// ListAssignments returns a user's per-server MCP access assignments in a scope.
func (r *MCPResource) ListAssignments(ctx context.Context, scopeID, subjectID string) ([]McpUserAssignment, error) {
	type output struct {
		Assignments []McpUserAssignment `json:"assignments"`
	}
	var out output
	path := fmt.Sprintf("/v1/scope/%s/mcp-assignment/%s", scopeID, url.PathEscape(subjectID))
	if err := r.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Assignments, nil
}

// SetAssignment grants (or updates) a user's access to an MCP server with a max trust tier.
// Optional resourceConstraints restricts which argument values the user may pass (ABAC fallback).
func (r *MCPResource) SetAssignment(
	ctx context.Context,
	scopeID, subjectID, serverKey string,
	maxTier TrustTier,
	resourceConstraints map[string][]string,
) (*McpUserAssignment, error) {
	var result McpUserAssignment
	body := map[string]any{"maxTier": maxTier}
	if len(resourceConstraints) > 0 {
		body["resourceConstraints"] = resourceConstraints
	}
	path := fmt.Sprintf("/v1/scope/%s/mcp-assignment/%s/%s", scopeID, url.PathEscape(subjectID), url.PathEscape(serverKey))
	if err := r.client.put(ctx, path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAssignment revokes a user's access to an MCP server.
func (r *MCPResource) DeleteAssignment(ctx context.Context, scopeID, subjectID, serverKey string) error {
	path := fmt.Sprintf("/v1/scope/%s/mcp-assignment/%s/%s", scopeID, url.PathEscape(subjectID), url.PathEscape(serverKey))
	return r.client.delete(ctx, path)
}

// ProvisionForOrg provisions an MCP server for every active team member at a capped trust tier
// (Enterprise-Managed Authorization — zero-touch, no per-user OAuth). Returns the count.
func (r *MCPResource) ProvisionForOrg(ctx context.Context, scopeID, serverKey string, maxTier TrustTier) (int, error) {
	type output struct {
		Provisioned int `json:"provisioned"`
	}
	var out output
	path := fmt.Sprintf("/v1/scope/%s/mcp-assignment/provision", scopeID)
	if err := r.client.post(ctx, path, map[string]any{"serverKey": serverKey, "maxTier": maxTier}, &out); err != nil {
		return 0, err
	}
	return out.Provisioned, nil
}
