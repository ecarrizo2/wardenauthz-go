package wardenauth

import (
	"context"
	"fmt"
)

type AgentResource struct{ client *Client }

func (r *AgentResource) Identify(ctx context.Context, scopeID string, input AgentIdentifyInput) (*AgentIdentifyResult, error) {
	var result AgentIdentifyResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/agent/identify", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AgentResource) Check(ctx context.Context, scopeID string, input AgentCheckInput) (*AgentCheckResult, error) {
	var result AgentCheckResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/agent/check", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
