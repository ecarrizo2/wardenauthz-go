package wardenauth

import (
	"context"
	"fmt"
	"net/url"
)

type SessionTokensResource struct{ client *Client }

func (r *SessionTokensResource) Mint(ctx context.Context, input *MintSessionTokenInput) (*MintSessionTokenResult, error) {
	var result MintSessionTokenResult
	if err := r.client.post(ctx, "/v1/session-token/mint", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SessionTokensResource) MintWithIntent(ctx context.Context, input MintIntentSessionTokenInput) (*MintIntentSessionTokenResult, error) {
	var result MintIntentSessionTokenResult
	if err := r.client.post(ctx, "/v1/session-token/intent-mint", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SessionTokensResource) VerifyIntentCall(ctx context.Context, input VerifyIntentCallInput) (*VerifyIntentCallResult, error) {
	var result VerifyIntentCallResult
	if err := r.client.post(ctx, "/v1/session-token/verify-intent-call", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type RevokeSessionTokenResult struct {
	JTI            string `json:"jti"`
	Revoked        bool   `json:"revoked"`
	AlreadyRevoked bool   `json:"alreadyRevoked,omitempty"`
}

func (r *SessionTokensResource) Revoke(ctx context.Context, jti string) (*RevokeSessionTokenResult, error) {
	var result RevokeSessionTokenResult
	if err := r.client.do(ctx, "DELETE", fmt.Sprintf("/v1/session-token/%s", jti), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type SodConstraintsResource struct{ client *Client }

func (r *SodConstraintsResource) Create(ctx context.Context, scopeID string, input CreateSodConstraintInput) (*SodConstraintItem, error) {
	var result SodConstraintItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/sod-constraint", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SodConstraintsResource) List(ctx context.Context, scopeID string) ([]SodConstraintItem, error) {
	type output struct{ Items []SodConstraintItem `json:"items"` }
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/sod-constraint", scopeID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *SodConstraintsResource) Delete(ctx context.Context, scopeID, constraintID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/sod-constraint/%s", scopeID, constraintID))
}

type TeamMembersResource struct{ client *Client }

func (r *TeamMembersResource) Add(ctx context.Context, scopeID string, input AddTeamMemberInput) (*TeamMemberItem, error) {
	var result TeamMemberItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/team-members", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *TeamMembersResource) List(ctx context.Context, scopeID string) ([]TeamMemberItem, error) {
	type output struct{ Items []TeamMemberItem `json:"items"` }
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/team-members", scopeID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *TeamMembersResource) Remove(ctx context.Context, scopeID, subjectID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/team-members/%s", scopeID, subjectID))
}

type ResourceTypesResource struct{ client *Client }

func (r *ResourceTypesResource) Create(ctx context.Context, scopeID string, input CreateResourceTypeInput) (*ResourceTypeItem, error) {
	var result ResourceTypeItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/resource-type", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResourceTypesResource) List(ctx context.Context, scopeID string, opts *struct{ Limit, NextToken string }) (*PaginatedResult[ResourceTypeItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
	}
	var result PaginatedResult[ResourceTypeItem]
	if err := r.client.get(ctx, buildQueryPath(fmt.Sprintf("/v1/scope/%s/resource-type", scopeID), params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResourceTypesResource) Delete(ctx context.Context, scopeID, typeID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/resource-type/%s", scopeID, typeID))
}

type ConsentResource struct{ client *Client }

func (r *ConsentResource) GetServers(ctx context.Context, scopeID string) ([]McpConsentServerSummary, error) {
	type output struct{ Servers []McpConsentServerSummary `json:"servers"` }
	var out output
	path := "/v1/mcp/consent/servers"
	if scopeID != "" {
		path += "?scopeId=" + url.QueryEscape(scopeID)
	}
	if err := r.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Servers, nil
}

func (r *ConsentResource) GetContext(ctx context.Context, serverKey, scopeID string) (*McpConsentContext, error) {
	var result McpConsentContext
	path := fmt.Sprintf("/v1/mcp/consent/context?serverKey=%s", url.QueryEscape(serverKey))
	if scopeID != "" {
		path += "&scopeId=" + url.QueryEscape(scopeID)
	}
	if err := r.client.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ConsentResource) GetPortalContext(ctx context.Context, scopeID string) ([]McpConsentContext, error) {
	type output struct{ Servers []McpConsentContext `json:"servers"` }
	var out output
	path := "/v1/mcp/consent/context"
	if scopeID != "" {
		path += "?scopeId=" + url.QueryEscape(scopeID)
	}
	if err := r.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Servers, nil
}

func (r *ConsentResource) Grant(ctx context.Context, input McpConsentGrantBody) (*McpConsentGrantResult, error) {
	var result McpConsentGrantResult
	if err := r.client.post(ctx, "/v1/mcp/consent/grant", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ConsentResource) Deny(ctx context.Context, authRequestID string) error {
	return r.client.post(ctx, "/v1/mcp/consent/deny", map[string]string{"authRequestId": authRequestID}, nil)
}

func (r *ConsentResource) GetRequestInfo(ctx context.Context, reqID string) (map[string]any, error) {
	var result map[string]any
	if err := r.client.get(ctx, fmt.Sprintf("/v1/mcp/consent/request?req=%s", url.QueryEscape(reqID)), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ConsentResource) ListGrants(ctx context.Context) ([]McpGrantSummary, error) {
	type output struct{ Grants []McpGrantSummary `json:"grants"` }
	var out output
	if err := r.client.get(ctx, "/v1/mcp/grants", &out); err != nil {
		return nil, err
	}
	return out.Grants, nil
}

func (r *ConsentResource) RevokeGrant(ctx context.Context, grantID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/mcp/grants/%s", url.PathEscape(grantID)))
}

func (r *ConsentResource) ListApprovals(ctx context.Context) ([]McpApprovalSummary, error) {
	type output struct{ Approvals []McpApprovalSummary `json:"approvals"` }
	var out output
	if err := r.client.get(ctx, "/v1/mcp/approvals", &out); err != nil {
		return nil, err
	}
	return out.Approvals, nil
}

func (r *ConsentResource) ListApprovalHistory(ctx context.Context) ([]McpApprovalHistoryItem, error) {
	type output struct{ Approvals []McpApprovalHistoryItem `json:"approvals"` }
	var out output
	if err := r.client.get(ctx, "/v1/mcp/approvals/history", &out); err != nil {
		return nil, err
	}
	return out.Approvals, nil
}

func (r *ConsentResource) ApproveRequest(ctx context.Context, id string, reason *string) error {
	body := map[string]any{}
	if reason != nil {
		body["reason"] = *reason
	}
	return r.client.post(ctx, fmt.Sprintf("/v1/mcp/approvals/%s/approve", url.PathEscape(id)), body, nil)
}

func (r *ConsentResource) DenyRequest(ctx context.Context, id string, reason *string) error {
	body := map[string]any{}
	if reason != nil {
		body["reason"] = *reason
	}
	return r.client.post(ctx, fmt.Sprintf("/v1/mcp/approvals/%s/deny", url.PathEscape(id)), body, nil)
}

func (r *ConsentResource) GetPushPublicKey(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	if err := r.client.get(ctx, "/v1/mcp/push/public-key", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ConsentResource) SubscribePush(ctx context.Context, endpoint string, p256dh, auth string) error {
	body := map[string]any{
		"endpoint": endpoint,
		"keys":     map[string]string{"p256dh": p256dh, "auth": auth},
	}
	return r.client.post(ctx, "/v1/mcp/push/subscribe", body, nil)
}

func (r *ConsentResource) UnsubscribePush(ctx context.Context, endpoint string) error {
	return r.client.post(ctx, "/v1/mcp/push/unsubscribe", map[string]string{"endpoint": endpoint}, nil)
}

func (r *ConsentResource) GetVelocityConfig(ctx context.Context) (*McpVelocityConfig, error) {
	var result McpVelocityConfig
	if err := r.client.get(ctx, "/v1/mcp/velocity-config", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ConsentResource) UpdateVelocityConfig(ctx context.Context, input map[string]any) error {
	return r.client.patch(ctx, "/v1/mcp/velocity-config", input, nil)
}
