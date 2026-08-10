package wardenauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientConfig struct {
	APIURL     string
	APIKey     string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

type Client struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	maxRetries     int
	retryDelay     time.Duration
	Scopes         *ScopesResource
	Permissions    *PermissionsResource
	Roles          *RolesResource
	AccessPolicies *AccessPoliciesResource
	APIKeys        *APIKeysResource
	Webhooks       *WebhooksResource
	Access         *AccessResource
	Audit          *AuditResource
	SessionTokens  *SessionTokensResource
	SodConstraints *SodConstraintsResource
	TeamMembers    *TeamMembersResource
	ResourceTypes  *ResourceTypesResource
	Tuples         *TuplesResource
	Integrations   *IntegrationsResource
	MCP            *MCPResource
	Consent        *ConsentResource
	Agent          *AgentResource
}

func NewClient(cfg ClientConfig) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	retryDelay := cfg.RetryDelay
	if retryDelay == 0 {
		retryDelay = defaultRetryDelay
	}

	c := &Client{
		baseURL: strings.TrimRight(cfg.APIURL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: timeout,
			// Tuned transport: the default http.Transport caps MaxIdleConnsPerHost at 2,
			// which forces a fresh TLS handshake on most concurrent calls to this single-host
			// API. Raising idle-conn reuse keeps connections warm and removes that per-call
			// handshake latency. HTTP/2 multiplexing is also enabled.
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   64,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
	c.Scopes = &ScopesResource{c}
	c.Permissions = &PermissionsResource{c}
	c.Roles = &RolesResource{c}
	c.AccessPolicies = &AccessPoliciesResource{c}
	c.APIKeys = &APIKeysResource{c}
	c.Webhooks = &WebhooksResource{c}
	c.Access = &AccessResource{c}
	c.Audit = &AuditResource{c}
	c.SessionTokens = &SessionTokensResource{c}
	c.SodConstraints = &SodConstraintsResource{c}
	c.TeamMembers = &TeamMembersResource{c}
	c.ResourceTypes = &ResourceTypesResource{c}
	c.Tuples = &TuplesResource{c}
	c.Integrations = &IntegrationsResource{c}
	c.MCP = &MCPResource{c}
	c.Consent = &ConsentResource{c}
	c.Agent = &AgentResource{c}
	return c
}

type APIError struct {
	StatusCode int
	Body       json.RawMessage
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("warden-auth api error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("warden-auth api error %d: %s", e.StatusCode, string(e.Body))
}

func (c *Client) do(ctx context.Context, method, path string, body any, result any) error {
	return c.doWithRetry(ctx, method, path, body, result)
}

func (c *Client) doRaw(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       respBody,
			Message:    resp.Status,
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) post(ctx context.Context, path string, body any, result any) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *Client) patch(ctx context.Context, path string, body any, result any) error {
	return c.do(ctx, "PATCH", path, body, result)
}

func (c *Client) put(ctx context.Context, path string, body any, result any) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

func (c *Client) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) getRaw(ctx context.Context, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Body:       respBody,
			Message:    resp.Status,
		}
	}

	return string(respBody), nil
}

func buildQueryPath(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}
	values := url.Values{}
	for k, v := range params {
		if v != "" {
			values.Set(k, v)
		}
	}
	return path + "?" + values.Encode()
}
