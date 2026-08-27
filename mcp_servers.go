package wardenauth

import (
	"context"
	"fmt"
)

type McpServersResource struct{ client *Client }

func (r *McpServersResource) Create(ctx context.Context, scopeID string, input CreateMcpServerInput) (*McpServer, error) {
	var result McpServer
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/mcp-server", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *McpServersResource) List(ctx context.Context, scopeID string) ([]McpServer, error) {
	type output struct {
		Items []McpServer `json:"items"`
	}
	var out output
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/mcp-server", scopeID), &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (r *McpServersResource) GetByID(ctx context.Context, scopeID, id string) (*McpServer, error) {
	var result McpServer
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/mcp-server/%s", scopeID, id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *McpServersResource) Update(ctx context.Context, scopeID, id string, input UpdateMcpServerInput) (*McpServer, error) {
	var result McpServer
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/mcp-server/%s", scopeID, id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *McpServersResource) Delete(ctx context.Context, scopeID, id string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/mcp-server/%s", scopeID, id))
}
