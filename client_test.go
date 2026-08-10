package wardenauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// ─── HTTP Method Helpers ────────────────────────────────────────────────────

func (c *Client) getPublic(ctx context.Context, path string, result any) error {
	return c.get(ctx, path, result)
}

func (c *Client) postPublic(ctx context.Context, path string, body any, result any) error {
	return c.post(ctx, path, body, result)
}

func (c *Client) patchPublic(ctx context.Context, path string, body any, result any) error {
	return c.patch(ctx, path, body, result)
}

func (c *Client) putPublic(ctx context.Context, path string, body any, result any) error {
	return c.put(ctx, path, body, result)
}

func (c *Client) deletePublic(ctx context.Context, path string) error {
	return c.delete(ctx, path)
}

func (c *Client) doRawPublic(ctx context.Context, method, path string, body any, result any) error {
	return c.doRaw(ctx, method, path, body, result)
}

func (c *Client) getRawPublic(ctx context.Context, path string) (string, error) {
	return c.getRaw(ctx, path)
}

// ─── NewClient ───────────────────────────────────────────────────────────────

func TestNewClient_DefaultConfig(t *testing.T) {
	client := NewClient(ClientConfig{
		APIURL: "https://api.example.com",
		APIKey: "sk_test",
	})

	if client.baseURL != "https://api.example.com" {
		t.Errorf("expected baseURL %q, got %q", "https://api.example.com", client.baseURL)
	}
	if client.apiKey != "sk_test" {
		t.Errorf("expected apiKey %q, got %q", "sk_test", client.apiKey)
	}
	if client.maxRetries != defaultMaxRetries {
		t.Errorf("expected maxRetries %d, got %d", defaultMaxRetries, client.maxRetries)
	}
	if client.retryDelay != defaultRetryDelay {
		t.Errorf("expected retryDelay %v, got %v", defaultRetryDelay, client.retryDelay)
	}
	if client.httpClient == nil {
		t.Fatal("httpClient should not be nil")
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout %v, got %v", 30*time.Second, client.httpClient.Timeout)
	}

	if client.Scopes == nil {
		t.Error("Scopes resource should be initialized")
	}
	if client.Permissions == nil {
		t.Error("Permissions resource should be initialized")
	}
}

func TestNewClient_BaseURLTrailingSlash(t *testing.T) {
	client := NewClient(ClientConfig{
		APIURL: "https://api.example.com/",
		APIKey: "sk_test",
	})

	if client.baseURL != "https://api.example.com" {
		t.Errorf("expected trailing slash to be trimmed, got %q", client.baseURL)
	}
}

func TestNewClient_CustomConfig(t *testing.T) {
	client := NewClient(ClientConfig{
		APIURL:     "https://api.example.com",
		APIKey:     "sk_test",
		Timeout:    5 * time.Second,
		MaxRetries: 5,
		RetryDelay: 1 * time.Second,
	})

	if client.maxRetries != 5 {
		t.Errorf("expected maxRetries 5, got %d", client.maxRetries)
	}
	if client.retryDelay != 1*time.Second {
		t.Errorf("expected retryDelay 1s, got %v", client.retryDelay)
	}
	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", client.httpClient.Timeout)
	}
}

func TestNewClient_ZeroRetriesMeansDefault(t *testing.T) {
	client := NewClient(ClientConfig{
		APIURL:     "https://api.example.com",
		APIKey:     "sk_test",
		MaxRetries: 0,
	})

	if client.maxRetries != defaultMaxRetries {
		t.Errorf("expected default retries when 0, got %d", client.maxRetries)
	}
}

// ─── GET ─────────────────────────────────────────────────────────────────────

func TestClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("expected x-api-key header, got %q", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Authorization header, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"scope1","name":"Test Scope","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.getPublic(context.Background(), "/v1/scope", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "scope1" {
		t.Errorf("expected ID scope1, got %q", result.ID)
	}
	if result.Name != "Test Scope" {
		t.Errorf("expected Name 'Test Scope', got %q", result.Name)
	}
}

func TestClient_Get_PaginatedResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"p1","name":"perm1"}],"nextToken":"next-page"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	type simpleItem struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var result PaginatedResult[simpleItem]
	err := client.getPublic(context.Background(), "/v1/scope/s1/permission", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].ID != "p1" {
		t.Errorf("expected p1, got %q", result.Items[0].ID)
	}
	if result.NextToken != "next-page" {
		t.Errorf("expected nextToken 'next-page', got %q", result.NextToken)
	}
}

func TestClient_Get_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	err := client.getPublic(context.Background(), "/v1/scope/nonexistent", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

// ─── POST ────────────────────────────────────────────────────────────────────

func TestClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}

		var body CreateScopeInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.ID != "scope-new" {
			t.Errorf("expected ID scope-new, got %q", body.ID)
		}
		if body.Name != "New Scope" {
			t.Errorf("expected Name 'New Scope', got %q", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"scope-new","name":"New Scope","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.postPublic(context.Background(), "/v1/scope", CreateScopeInput{
		ID:   "scope-new",
		Name: "New Scope",
		Type: "workspace",
	}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "scope-new" {
		t.Errorf("expected ID scope-new, got %q", result.ID)
	}
}

func TestClient_Post_NoBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result map[string]any
	err := client.postPublic(context.Background(), "/v1/action", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["message"] != "ok" {
		t.Errorf("expected message 'ok', got %v", result["message"])
	}
}

func TestClient_Post_NonRetryableError(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	err := client.postPublic(context.Background(), "/v1/scope", CreateScopeInput{
		ID:   "bad",
		Name: "Bad",
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}

	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 call for non-retryable error, got %d", callCount)
	}
}

// ─── PATCH ───────────────────────────────────────────────────────────────────

func TestClient_Patch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		name := "Updated"
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"id":"s1","name":"%s","type":"workspace"}`, name)))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.patchPublic(context.Background(), "/v1/scope/s1", UpdateScopeInput{
		Name: stringPtr("Updated"),
	}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated" {
		t.Errorf("expected Name 'Updated', got %q", result.Name)
	}
}

// ─── PUT ─────────────────────────────────────────────────────────────────────

func TestClient_Put_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["maxTier"] != "low" {
			t.Errorf("expected maxTier low, got %v", body["maxTier"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subjectId":"user1","serverKey":"server1","maxTier":"low","assignedBy":"admin","createdAt":123,"updatedAt":456}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result McpUserAssignment
	err := client.putPublic(context.Background(), "/v1/scope/s1/mcp-assignment/u1/s1", map[string]any{"maxTier": "low"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MaxTier != "low" {
		t.Errorf("expected maxTier low, got %q", result.MaxTier)
	}
}

// ─── DELETE ──────────────────────────────────────────────────────────────────

func TestClient_Delete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	err := client.deletePublic(context.Background(), "/v1/scope/s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":"cannot delete scope with children"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	err := client.deletePublic(context.Background(), "/v1/scope/s1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("expected status 409, got %d", apiErr.StatusCode)
	}
}

// ─── doRaw ───────────────────────────────────────────────────────────────────

func TestClient_DoRaw_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body CreateScopeInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.ID != "raw-scope" {
			t.Errorf("expected ID raw-scope, got %q", body.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"raw-scope","name":"Raw Test","type":"organization"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.doRawPublic(context.Background(), http.MethodPost, "/v1/scope", CreateScopeInput{
		ID:   "raw-scope",
		Name: "Raw Test",
		Type: "organization",
	}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "raw-scope" {
		t.Errorf("expected ID raw-scope, got %q", result.ID)
	}
}

func TestClient_DoRaw_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid api key"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "bad-key"})

	err := client.doRawPublic(context.Background(), http.MethodGet, "/v1/scope", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
}

// ─── getRaw ──────────────────────────────────────────────────────────────────

func TestClient_GetRaw_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte("id,name,type\n1,Scope1,workspace\n"))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	body, err := client.getRawPublic(context.Background(), "/v1/audit/export")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, "id,name,type") {
		t.Errorf("expected CSV content, got %q", body)
	}
}

func TestClient_GetRaw_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"forbidden"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	_, err := client.getRawPublic(context.Background(), "/v1/audit/export")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ─── buildQueryPath ──────────────────────────────────────────────────────────

func TestBuildQueryPath_Empty(t *testing.T) {
	result := buildQueryPath("/v1/scope", nil)
	if result != "/v1/scope" {
		t.Errorf("expected '/v1/scope', got %q", result)
	}

	result = buildQueryPath("/v1/scope", map[string]string{})
	if result != "/v1/scope" {
		t.Errorf("expected '/v1/scope', got %q", result)
	}
}

func TestBuildQueryPath_SingleParam(t *testing.T) {
	result := buildQueryPath("/v1/scope", map[string]string{
		"limit": "50",
	})
	if result != "/v1/scope?limit=50" {
		t.Errorf("expected '/v1/scope?limit=50', got %q", result)
	}
}

func TestBuildQueryPath_MultipleParams(t *testing.T) {
	result := buildQueryPath("/v1/scope", map[string]string{
		"limit":     "50",
		"nextToken": "abc123",
		"type":      "workspace",
	})
	if !strings.Contains(result, "limit=50") {
		t.Errorf("expected param limit=50 in %q", result)
	}
	if !strings.Contains(result, "nextToken=abc123") {
		t.Errorf("expected param nextToken=abc123 in %q", result)
	}
	if !strings.Contains(result, "type=workspace") {
		t.Errorf("expected param type=workspace in %q", result)
	}
	if !strings.HasPrefix(result, "/v1/scope?") {
		t.Errorf("expected path prefix, got %q", result)
	}
}

func TestBuildQueryPath_SkipsEmptyValues(t *testing.T) {
	result := buildQueryPath("/v1/scope", map[string]string{
		"limit":     "50",
		"nextToken": "",
		"type":      "",
	})
	if strings.Contains(result, "nextToken") {
		t.Errorf("expected empty nextToken to be skipped, got %q", result)
	}
	if strings.Contains(result, "type") {
		t.Errorf("expected empty type to be skipped, got %q", result)
	}
	if !strings.Contains(result, "limit=50") {
		t.Errorf("expected limit=50 in %q", result)
	}
}

func TestBuildQueryPath_SpecialCharacters(t *testing.T) {
	result := buildQueryPath("/v1/scope", map[string]string{
		"name": "test&scope=1",
	})
	if !strings.Contains(result, "name=test%26scope%3D1") {
		t.Errorf("expected URL-encoded value, got %q", result)
	}
}

// ─── APIError ────────────────────────────────────────────────────────────────

func TestAPIError_Error_WithMessage(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "not found",
		Body:       json.RawMessage(`{"error":"not found"}`),
	}
	expected := "warden-auth api error 404: not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestAPIError_Error_WithBody(t *testing.T) {
	err := &APIError{
		StatusCode: 500,
		Message:    "",
		Body:       json.RawMessage(`{"error":"internal"}`),
	}
	expected := `warden-auth api error 500: {"error":"internal"}`
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestAPIError_Error_EmptyBody(t *testing.T) {
	err := &APIError{
		StatusCode: 429,
		Message:    "",
		Body:       json.RawMessage{},
	}
	expected := "warden-auth api error 429: "
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

// ─── isRetryable ─────────────────────────────────────────────────────────────

func TestIsRetryable_RetryableStatusCodes(t *testing.T) {
	retryableStatuses := []int{429, 500, 502, 503, 504, 599}
	for _, code := range retryableStatuses {
		if !isRetryable(code) {
			t.Errorf("status %d should be retryable", code)
		}
	}
}

func TestIsRetryable_NonRetryableStatusCodes(t *testing.T) {
	nonRetryable := []int{200, 201, 204, 301, 302, 400, 401, 403, 404, 409, 422, 428}
	for _, code := range nonRetryable {
		if isRetryable(code) {
			t.Errorf("status %d should NOT be retryable", code)
		}
	}
}

// ─── Retry Behavior ──────────────────────────────────────────────────────────

func TestClient_Retry_On429(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"s1","name":"Scope","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 3,
	})

	var result ScopeItem
	err := client.getPublic(context.Background(), "/v1/scope/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "s1" {
		t.Errorf("expected ID s1, got %q", result.ID)
	}
	if atomic.LoadInt32(&callCount) != 3 {
		t.Errorf("expected 3 calls (2 retries + 1 success), got %d", atomic.LoadInt32(&callCount))
	}
}

func TestClient_Retry_On5xx(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error":"temporary"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"ok","name":"Final","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 3,
	})

	var result ScopeItem
	err := client.getPublic(context.Background(), "/v1/scope/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Final" {
		t.Errorf("expected Name 'Final', got %q", result.Name)
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 calls, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestClient_Retry_Exhausted(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"bad gateway"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 2,
	})

	err := client.getPublic(context.Background(), "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error after retry exhaustion, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502, got %d", apiErr.StatusCode)
	}
	if atomic.LoadInt32(&callCount) != 3 {
		t.Errorf("expected 3 calls (initial + 2 retries), got %d", atomic.LoadInt32(&callCount))
	}
}

func TestClient_NoRetry_On400(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"validation failed"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 3,
	})

	err := client.getPublic(context.Background(), "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected exactly 1 call (no retries for 4xx), got %d", atomic.LoadInt32(&callCount))
	}
}

func TestClient_NoRetry_On401(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 3,
	})

	err := client.getPublic(context.Background(), "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 call for 401, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestClient_NoRetry_On403(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"forbidden"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 1 * time.Millisecond,
		MaxRetries: 3,
	})

	err := client.getPublic(context.Background(), "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 call for 403, got %d", atomic.LoadInt32(&callCount))
	}
}

// ─── Context Cancellation ────────────────────────────────────────────────────

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL: server.URL,
		APIKey: "test-key",
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.getPublic(ctx, "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}

func TestClient_ContextDeadline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:  server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := client.getPublic(ctx, "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected deadline exceeded error, got nil")
	}
}

func TestClient_ContextCancelledDuringRetry(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"rate limited"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIURL:     server.URL,
		APIKey:     "test-key",
		RetryDelay: 500 * time.Millisecond,
		MaxRetries: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := client.getPublic(ctx, "/v1/scope/test", nil)
	if err == nil {
		t.Fatal("expected error from cancelled context during retry, got nil")
	}
}

// ─── Integration-style Tests (via T.Run) ─────────────────────────────────────

func TestClient_MethodNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	if err := client.getPublic(context.Background(), "/v1/scope/s1", &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "ok" {
		t.Errorf("expected ID ok, got %q", result.ID)
	}
}

func TestClient_ConcurrentRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"concurrent","name":"Test","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	errCh := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			var result ScopeItem
			errCh <- client.getPublic(context.Background(), "/v1/scope/test", &result)
		}()
	}

	for i := 0; i < 10; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("concurrent request %d failed: %v", i, err)
		}
	}
}

func TestClient_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.getPublic(context.Background(), "/v1/scope/test", &result)
	if err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("expected unmarshal error, got %q", err.Error())
	}
}

func TestClient_EmptyResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})

	var result ScopeItem
	err := client.getPublic(context.Background(), "/v1/scope/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_CustomHeadersOnRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "custom-key" {
			t.Errorf("expected x-api-key 'custom-key', got %q", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("Authorization") != "Bearer custom-key" {
			t.Errorf("expected Authorization 'Bearer custom-key', got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "custom-key"})

	err := client.deletePublic(context.Background(), "/v1/scope/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ─── Paginate ────────────────────────────────────────────────────────────────

func TestClient_Paginate_SinglePage(t *testing.T) {
	client := NewClient(ClientConfig{APIURL: "https://api.example.com", APIKey: "test-key"})

	callCount := 0
	err := client.Paginate(context.Background(), func(ctx context.Context, nextToken *string) (*string, error) {
		callCount++
		if nextToken != nil {
			t.Errorf("expected nil nextToken on first call, got %v", *nextToken)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestClient_Paginate_MultiplePages(t *testing.T) {
	client := NewClient(ClientConfig{APIURL: "https://api.example.com", APIKey: "test-key"})

	pageTokens := []string{"page2", "page3", ""}
	pageIdx := 0
	callCount := 0
	err := client.Paginate(context.Background(), func(ctx context.Context, nextToken *string) (*string, error) {
		callCount++
		result := pageTokens[pageIdx]
		pageIdx++
		if result == "" {
			return nil, nil
		}
		return stringPtr(result), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestClient_Paginate_Error(t *testing.T) {
	client := NewClient(ClientConfig{APIURL: "https://api.example.com", APIKey: "test-key"})

	err := client.Paginate(context.Background(), func(ctx context.Context, nextToken *string) (*string, error) {
		return nil, fmt.Errorf("pagination error")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "pagination error" {
		t.Errorf("expected 'pagination error', got %q", err.Error())
	}
}

func TestClient_Paginate_ContextCancelled(t *testing.T) {
	client := NewClient(ClientConfig{APIURL: "https://api.example.com", APIKey: "test-key"})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Paginate(ctx, func(ctx context.Context, nextToken *string) (*string, error) {
		t.Error("callback should not be called when context is cancelled")
		return nil, nil
	})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestClient_Paginate_ContextCancelledMidIteration(t *testing.T) {
	client := NewClient(ClientConfig{APIURL: "https://api.example.com", APIKey: "test-key"})

	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	err := client.Paginate(ctx, func(ctxArg context.Context, nextToken *string) (*string, error) {
		callCount++
		if callCount == 2 {
			cancel()
		}
		return stringPtr("next"), nil
	})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if callCount < 2 {
		t.Errorf("expected at least 2 callbacks, got %d", callCount)
	}
}

// ─── Helper Functions ────────────────────────────────────────────────────────

func TestStringPtr(t *testing.T) {
	s := "hello"
	p := stringPtr(s)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != "hello" {
		t.Errorf("expected 'hello', got %q", *p)
	}
}

func TestIntPtr(t *testing.T) {
	i := 42
	p := intPtr(i)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != 42 {
		t.Errorf("expected 42, got %d", *p)
	}
}

func TestBoolPtr(t *testing.T) {
	b := true
	p := boolPtr(b)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != true {
		t.Error("expected true")
	}

	b = false
	p = boolPtr(b)
	if *p != false {
		t.Error("expected false")
	}
}
