package wardenauth

import (
	"context"
	"fmt"
)

type WebhooksResource struct{ client *Client }

func (r *WebhooksResource) Create(ctx context.Context, scopeID string, input CreateWebhookEndpointInput) (*WebhookEndpointCreatedItem, error) {
	var result WebhookEndpointCreatedItem
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/webhooks", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *WebhooksResource) List(ctx context.Context, scopeID string, opts *struct{ Limit, NextToken string }) (*PaginatedResult[WebhookEndpointItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
	}
	var result PaginatedResult[WebhookEndpointItem]
	if err := r.client.get(ctx, buildQueryPath(fmt.Sprintf("/v1/scope/%s/webhooks", scopeID), params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *WebhooksResource) GetByID(ctx context.Context, scopeID, endpointID string) (*WebhookEndpointItem, error) {
	var result WebhookEndpointItem
	if err := r.client.get(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s", scopeID, endpointID), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *WebhooksResource) Update(ctx context.Context, scopeID, endpointID string, input UpdateWebhookEndpointInput) (*WebhookEndpointItem, error) {
	var result WebhookEndpointItem
	if err := r.client.patch(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s", scopeID, endpointID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *WebhooksResource) Delete(ctx context.Context, scopeID, endpointID string) error {
	return r.client.delete(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s", scopeID, endpointID))
}

func (r *WebhooksResource) RotateSecret(ctx context.Context, scopeID, endpointID string) (*WebhookRotateSecretResult, error) {
	var result WebhookRotateSecretResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s/rotate-secret", scopeID, endpointID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *WebhooksResource) Test(ctx context.Context, scopeID, endpointID string) (string, error) {
	type out struct{ Message string `json:"message"` }
	var output out
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s/test", scopeID, endpointID), nil, &output); err != nil {
		return "", err
	}
	return output.Message, nil
}

func (r *WebhooksResource) RetryDelivery(ctx context.Context, scopeID, endpointID, deliveryID string) (string, error) {
	type out struct{ Message string `json:"message"` }
	var output out
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/webhooks/%s/deliveries/%s/retry", scopeID, endpointID, deliveryID), nil, &output); err != nil {
		return "", err
	}
	return output.Message, nil
}

func (r *WebhooksResource) ListDeliveries(ctx context.Context, scopeID, endpointID string, opts *struct{ Limit, NextToken string }) (*PaginatedResult[WebhookDeliveryItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
	}
	var result PaginatedResult[WebhookDeliveryItem]
	if err := r.client.get(ctx, buildQueryPath(fmt.Sprintf("/v1/scope/%s/webhooks/%s/deliveries", scopeID, endpointID), params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
