package wardenauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ─── ScopesResource ───────────────────────────────────────────────────────────

func TestScopesResource_List_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope" {
			t.Errorf("expected /v1/scope, got %q", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("expected x-api-key 'test-key', got %q", r.Header.Get("x-api-key"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"s1","name":"Scope 1","type":"workspace"}],"nextToken":""}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].ID != "s1" {
		t.Errorf("expected ID s1, got %q", result.Items[0].ID)
	}
	if result.Items[0].Name != "Scope 1" {
		t.Errorf("expected Name 'Scope 1', got %q", result.Items[0].Name)
	}
}

func TestScopesResource_List_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope" {
			t.Errorf("expected /v1/scope, got %q", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %q", query.Get("limit"))
		}
		if query.Get("nextToken") != "abc" {
			t.Errorf("expected nextToken=abc, got %q", query.Get("nextToken"))
		}
		if query.Get("type") != "workspace" {
			t.Errorf("expected type=workspace, got %q", query.Get("type"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"s2","name":"Scope 2","type":"workspace"}],"nextToken":"next"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.List(context.Background(), &struct {
		Limit     string
		NextToken string
		Type      ScopeType
	}{Limit: "10", NextToken: "abc", Type: ScopeTypeWorkspace})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.NextToken != "next" {
		t.Errorf("expected nextToken 'next', got %q", result.NextToken)
	}
}

func TestScopesResource_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1" {
			t.Errorf("expected /v1/scope/scope1, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"scope1","name":"Scope One","type":"organization"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.GetByID(context.Background(), "scope1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "scope1" {
		t.Errorf("expected ID scope1, got %q", result.ID)
	}
	if result.Name != "Scope One" {
		t.Errorf("expected Name 'Scope One', got %q", result.Name)
	}
	if result.Type != "organization" {
		t.Errorf("expected Type 'organization', got %q", result.Type)
	}
}

func TestScopesResource_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope" {
			t.Errorf("expected /v1/scope, got %q", r.URL.Path)
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
		if body.Type != "organization" {
			t.Errorf("expected Type 'organization', got %q", body.Type)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"scope-new","name":"New Scope","type":"organization"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.Create(context.Background(), CreateScopeInput{
		ID:   "scope-new",
		Name: "New Scope",
		Type: "organization",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "scope-new" {
		t.Errorf("expected ID scope-new, got %q", result.ID)
	}
}

func TestScopesResource_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1" {
			t.Errorf("expected /v1/scope/scope1, got %q", r.URL.Path)
		}

		var body UpdateScopeInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.ID != "scope1" {
			t.Errorf("expected ID scope1, got %q", body.ID)
		}
		if body.Name == nil || *body.Name != "Updated" {
			t.Errorf("expected Name 'Updated', got %v", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"scope1","name":"Updated","type":"workspace"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.Update(context.Background(), "scope1", UpdateScopeInput{
		ID:   "scope1",
		Name: stringPtr("Updated"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated" {
		t.Errorf("expected Name 'Updated', got %q", result.Name)
	}
}

func TestScopesResource_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1" {
			t.Errorf("expected /v1/scope/scope1, got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Scopes.Delete(context.Background(), "scope1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScopesResource_ApplyManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/apply" {
			t.Errorf("expected /v1/scope/scope1/apply, got %q", r.URL.Path)
		}

		var body ApplyScopeManifestInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"scopeId":"scope1","dryRun":false,"manifestHash":"abc123"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.ApplyManifest(context.Background(), "scope1", ApplyScopeManifestInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ScopeId != "scope1" {
		t.Errorf("expected ScopeId scope1, got %q", result.ScopeId)
	}
	if result.ManifestHash != "abc123" {
		t.Errorf("expected ManifestHash abc123, got %q", result.ManifestHash)
	}
}

func TestScopesResource_Export(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/export" {
			t.Errorf("expected /v1/scope/scope1/export, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"permissions":[],"roles":[],"policies":[]}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.Export(context.Background(), "scope1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

// ─── PermissionsResource ─────────────────────────────────────────────────────

func TestPermissionsResource_List_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission" {
			t.Errorf("expected /v1/scope/scope1/permission, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"p1","scopeId":"scope1","resource":"doc","action":"read","effect":"allow","name":"Read Docs"}],"nextToken":""}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.List(context.Background(), "scope1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].ID != "p1" {
		t.Errorf("expected ID p1, got %q", result.Items[0].ID)
	}
	if result.Items[0].Resource != "doc" {
		t.Errorf("expected Resource doc, got %q", result.Items[0].Resource)
	}
}

func TestPermissionsResource_List_Paginated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("limit") != "20" {
			t.Errorf("expected limit=20, got %q", query.Get("limit"))
		}
		if query.Get("nextToken") != "page2" {
			t.Errorf("expected nextToken=page2, got %q", query.Get("nextToken"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"p2","scopeId":"scope1","resource":"doc","action":"write","effect":"allow","name":"Write Docs"}],"nextToken":"page3"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.List(context.Background(), "scope1", &struct{ Limit, NextToken string }{
		Limit: "20", NextToken: "page2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NextToken != "page3" {
		t.Errorf("expected nextToken page3, got %q", result.NextToken)
	}
}

func TestPermissionsResource_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/perm1" {
			t.Errorf("expected /v1/scope/scope1/permission/perm1, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"perm1","scopeId":"scope1","resource":"doc","action":"read","effect":"allow","name":"Read Docs"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.GetByID(context.Background(), "scope1", "perm1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "perm1" {
		t.Errorf("expected ID perm1, got %q", result.ID)
	}
}

func TestPermissionsResource_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission" {
			t.Errorf("expected /v1/scope/scope1/permission, got %q", r.URL.Path)
		}

		raw, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		json.Unmarshal(raw, &body)
		if body["id"] != "perm-new" {
			t.Errorf("expected id perm-new, got %v", body["id"])
		}
		if body["scopeId"] != "scope1" {
			t.Errorf("expected scopeId scope1, got %v", body["scopeId"])
		}
		if body["resource"] != "doc" {
			t.Errorf("expected resource doc, got %v", body["resource"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"perm-new","scopeId":"scope1","resource":"doc","action":"read","effect":"allow","name":"Read"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.Create(context.Background(), "scope1", CreatePermissionInput{
		ID:       "perm-new",
		ScopeId:  "scope1",
		Resource: "doc",
		Action:   "read",
		Effect:   "allow",
		Name:     "Read",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "perm-new" {
		t.Errorf("expected ID perm-new, got %q", result.ID)
	}
}

func TestPermissionsResource_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/perm1" {
			t.Errorf("expected /v1/scope/scope1/permission/perm1, got %q", r.URL.Path)
		}

		var body UpdatePermissionInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name == nil || *body.Name != "Updated Name" {
			t.Errorf("expected Name 'Updated Name', got %v", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"perm1","scopeId":"scope1","resource":"doc","action":"read","effect":"allow","name":"Updated Name"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.Update(context.Background(), "scope1", "perm1", UpdatePermissionInput{
		Name: stringPtr("Updated Name"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated Name" {
		t.Errorf("expected Name 'Updated Name', got %q", result.Name)
	}
}

func TestPermissionsResource_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/perm1" {
			t.Errorf("expected /v1/scope/scope1/permission/perm1, got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Permissions.Delete(context.Background(), "scope1", "perm1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionsResource_BulkCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/bulk" {
			t.Errorf("expected /v1/scope/scope1/permission/bulk, got %q", r.URL.Path)
		}

		var body []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body) != 2 {
			t.Fatalf("expected 2 permissions, got %d", len(body))
		}
		if body[0]["id"] != "bulk-p1" {
			t.Errorf("expected id bulk-p1, got %v", body[0]["id"])
		}
		if body[1]["id"] != "bulk-p2" {
			t.Errorf("expected id bulk-p2, got %v", body[1]["id"])
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"created":2}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Permissions.BulkCreate(context.Background(), "scope1", []CreatePermissionInput{
		{ID: "bulk-p1", ScopeId: "scope1", Resource: "doc", Action: "read", Effect: "allow", Name: "P1"},
		{ID: "bulk-p2", ScopeId: "scope1", Resource: "doc", Action: "write", Effect: "allow", Name: "P2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionsResource_BulkDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/bulk" {
			t.Errorf("expected /v1/scope/scope1/permission/bulk, got %q", r.URL.Path)
		}

		var body map[string][]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body["ids"]) != 2 {
			t.Fatalf("expected 2 ids, got %d", len(body["ids"]))
		}
		if body["ids"][0] != "p1" || body["ids"][1] != "p2" {
			t.Errorf("expected ids [p1 p2], got %v", body["ids"])
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Permissions.BulkDelete(context.Background(), "scope1", []string{"p1", "p2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionsResource_ImportCSV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/permission/import-csv" {
			t.Errorf("expected /v1/scope/scope1/permission/import-csv, got %q", r.URL.Path)
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["csv"] != "id,resource,action\ndoc,read,allow" {
			t.Errorf("expected csv content, got %q", body["csv"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"created":1,"skipped":0,"errors":[]}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.ImportCSV(context.Background(), "scope1", "id,resource,action\ndoc,read,allow")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Created != 1 {
		t.Errorf("expected Created 1, got %v", result.Created)
	}
}

// ─── RolesResource ────────────────────────────────────────────────────────────

func TestRolesResource_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/role" {
			t.Errorf("expected /v1/scope/scope1/role, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"r1","name":"Admin","description":"Admin role"}],"nextToken":""}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Roles.List(context.Background(), "scope1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].ID != "r1" {
		t.Errorf("expected ID r1, got %q", result.Items[0].ID)
	}
	if result.Items[0].Name != "Admin" {
		t.Errorf("expected Name Admin, got %q", result.Items[0].Name)
	}
}

func TestRolesResource_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/role/admin" {
			t.Errorf("expected /v1/scope/scope1/role/admin, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"admin","name":"Admin Role"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Roles.GetByID(context.Background(), "scope1", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "admin" {
		t.Errorf("expected ID admin, got %q", result.ID)
	}
}

func TestRolesResource_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/role" {
			t.Errorf("expected /v1/scope/scope1/role, got %q", r.URL.Path)
		}

		raw, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		json.Unmarshal(raw, &body)
		if body["id"] != "role-new" {
			t.Errorf("expected id role-new, got %v", body["id"])
		}
		if body["scopeId"] != "scope1" {
			t.Errorf("expected scopeId scope1, got %v", body["scopeId"])
		}
		if body["name"] != "New Role" {
			t.Errorf("expected name 'New Role', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"role-new","name":"New Role"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Roles.Create(context.Background(), "scope1", CreateRoleInput{
		ID:            "role-new",
		ScopeId:       "scope1",
		Name:          "New Role",
		PermissionIds: []string{"p1", "p2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "role-new" {
		t.Errorf("expected ID role-new, got %q", result.ID)
	}
}

func TestRolesResource_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/role/admin" {
			t.Errorf("expected /v1/scope/scope1/role/admin, got %q", r.URL.Path)
		}

		var body UpdateRoleInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name != "Updated Role" {
			t.Errorf("expected Name 'Updated Role', got %q", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"admin","name":"Updated Role"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Roles.Update(context.Background(), "scope1", "admin", UpdateRoleInput{
		Name: "Updated Role",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated Role" {
		t.Errorf("expected Name 'Updated Role', got %q", result.Name)
	}
}

func TestRolesResource_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope1/role/admin" {
			t.Errorf("expected /v1/scope/scope1/role/admin, got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Roles.Delete(context.Background(), "scope1", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRolesResource_Clone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/scope2/role/clone-from/admin-template" {
			t.Errorf("expected /v1/scope/scope2/role/clone-from/admin-template, got %q", r.URL.Path)
		}
		if r.Body == nil {
			t.Fatal("expected non-nil body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"admin-clone","name":"Admin Clone"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Roles.Clone(context.Background(), "scope2", "admin-template")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "admin-clone" {
		t.Errorf("expected ID admin-clone, got %q", result.ID)
	}
}

// ─── APIKeysResource ──────────────────────────────────────────────────────────

func TestApiKeysResource_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management" {
			t.Errorf("expected /v1/scope/org1/api-key/management, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"keyId":"key1","type":"management","maskedKey":"warden_***abc","createdAt":"2024-01-01","subjectId":"user1","name":"My Key"}]}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	items, err := client.APIKeys.List(context.Background(), "org1", "management")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].KeyId != "key1" {
		t.Errorf("expected KeyId key1, got %q", items[0].KeyId)
	}
}

func TestApiKeysResource_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management/key1" {
			t.Errorf("expected /v1/scope/org1/api-key/management/key1, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keyId":"key1","type":"management","maskedKey":"warden_***abc","createdAt":"2024-01-01","subjectId":"user1","name":"My Key"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.GetByID(context.Background(), "org1", "management", "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyId != "key1" {
		t.Errorf("expected KeyId key1, got %q", result.KeyId)
	}
	if result.Name != "My Key" {
		t.Errorf("expected Name 'My Key', got %q", result.Name)
	}
}

func TestApiKeysResource_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/application" {
			t.Errorf("expected /v1/scope/org1/api-key/application, got %q", r.URL.Path)
		}

		var body CreateApiKeyInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name != "App Key" {
			t.Errorf("expected Name 'App Key', got %q", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"keyId":"app-key-1","type":"application","maskedKey":"warden_***xyz","createdAt":"2024-01-01","subjectId":"user1","name":"App Key","rawKey":"wk_s3cret"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.Create(context.Background(), "org1", "application", CreateApiKeyInput{
		Name: "App Key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyId != "app-key-1" {
		t.Errorf("expected KeyId app-key-1, got %q", result.KeyId)
	}
	if result.RawKey != "wk_s3cret" {
		t.Errorf("expected RawKey wk_s3cret, got %q", result.RawKey)
	}
}

func TestApiKeysResource_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management/key1" {
			t.Errorf("expected /v1/scope/org1/api-key/management/key1, got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.APIKeys.Delete(context.Background(), "org1", "management", "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApiKeysResource_Rotate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management/key1/rotate" {
			t.Errorf("expected /v1/scope/org1/api-key/management/key1/rotate, got %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"keyId":"key1-rot","maskedKey":"warden_***new","rawKey":"wk_new","name":"My Key","oldKeyId":"key1","overlapExpiresAt":"2024-02-01"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.Rotate(context.Background(), "org1", "management", "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RawKey != "wk_new" {
		t.Errorf("expected RawKey wk_new, got %q", result.RawKey)
	}
	if result.OldKeyId != "key1" {
		t.Errorf("expected OldKeyId key1, got %q", result.OldKeyId)
	}
}

func TestApiKeysResource_UpdateAutoRotation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management/key1/auto-rotation" {
			t.Errorf("expected /v1/scope/org1/api-key/management/key1/auto-rotation, got %q", r.URL.Path)
		}

		var body UpdateApiKeyAutoRotationInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Enabled != true {
			t.Errorf("expected Enabled true, got %v", body.Enabled)
		}
		if body.IntervalDays == nil || *body.IntervalDays != 30 {
			t.Errorf("expected IntervalDays 30, got %v", body.IntervalDays)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keyId":"key1","type":"management","maskedKey":"warden_***abc","createdAt":"2024-01-01","subjectId":"user1","name":"My Key","autoRotationEnabled":true,"autoRotationIntervalDays":30}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.UpdateAutoRotation(context.Background(), "org1", "management", "key1", UpdateApiKeyAutoRotationInput{
		Enabled:      true,
		IntervalDays: float64Ptr(30),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AutoRotationEnabled == nil || *result.AutoRotationEnabled != true {
		t.Errorf("expected AutoRotationEnabled true")
	}
}

func TestApiKeysResource_RevealRotation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/scope/org1/api-key/management/key1/rotation/reveal" {
			t.Errorf("expected /v1/scope/org1/api-key/management/key1/rotation/reveal, got %q", r.URL.Path)
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["rotationRef"] != "rot-123" {
			t.Errorf("expected rotationRef rot-123, got %q", body["rotationRef"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keyId":"key1","apiKey":"wk_revealed"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.RevealRotation(context.Background(), "org1", "management", "key1", "rot-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyId != "key1" {
		t.Errorf("expected KeyId key1, got %q", result.KeyId)
	}
	if result.ApiKey != "wk_revealed" {
		t.Errorf("expected ApiKey wk_revealed, got %q", result.ApiKey)
	}
}

// ─── Error Case Tests ─────────────────────────────────────────────────────────

func TestScopesResource_GetByID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"scope not found"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	_, err := client.Scopes.GetByID(context.Background(), "nonexistent")
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

func TestPermissionsResource_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":"permission in use"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	err := client.Permissions.Delete(context.Background(), "scope1", "perm1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("expected status 409, got %d", apiErr.StatusCode)
	}
}

// ─── Headers Verification ─────────────────────────────────────────────────────

func TestScopesResource_Delete_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "my-api-key" {
			t.Errorf("expected x-api-key 'my-api-key', got %q", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("Authorization") != "Bearer my-api-key" {
			t.Errorf("expected Authorization 'Bearer my-api-key', got %q", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "my-api-key"})
	err := client.Scopes.Delete(context.Background(), "scope1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ─── Response Parsing Verification ────────────────────────────────────────────

func TestScopesResource_List_ParsesAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"s1","name":"Scope 1","type":"organization","parentId":"parent1","description":"A scope","allowRoleInheritance":true,"accessScopeId":"acc1"}],"nextToken":"next"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Scopes.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	item := result.Items[0]
	if item.ID != "s1" {
		t.Errorf("expected ID s1, got %q", item.ID)
	}
	if item.ParentId == nil || *item.ParentId != "parent1" {
		t.Errorf("expected ParentId parent1, got %v", item.ParentId)
	}
	if item.Description == nil || *item.Description != "A scope" {
		t.Errorf("expected Description 'A scope', got %v", item.Description)
	}
	if item.AllowRoleInheritance == nil || *item.AllowRoleInheritance != true {
		t.Errorf("expected AllowRoleInheritance true, got %v", item.AllowRoleInheritance)
	}
	if item.AccessScopeId == nil || *item.AccessScopeId != "acc1" {
		t.Errorf("expected AccessScopeId acc1, got %v", item.AccessScopeId)
	}
	if result.NextToken != "next" {
		t.Errorf("expected NextToken next, got %q", result.NextToken)
	}
}

func TestPermissionsResource_List_ParsesAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[{"id":"p1","scopeId":"s1","resource":"doc","action":"read","effect":"deny","name":"Deny reads","description":"Blocks reading"}],"nextToken":""}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.Permissions.List(context.Background(), "scope1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item := result.Items[0]
	if item.ScopeId != "s1" {
		t.Errorf("expected ScopeId s1, got %q", item.ScopeId)
	}
	if item.Effect != "deny" {
		t.Errorf("expected Effect deny, got %q", item.Effect)
	}
	if item.Description == nil || *item.Description != "Blocks reading" {
		t.Errorf("expected Description 'Blocks reading', got %v", item.Description)
	}
}

func TestApiKeysResource_Create_ParsesAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"keyId":"k1","type":"management","maskedKey":"wk_***","createdAt":"2024-01-01","subjectId":"user1","name":"Key","scopeId":"org1","expiresAt":"2025-01-01","lastUsedAt":"2024-06-01","attributes":{"env":"prod"},"autoRotationEnabled":true,"autoRotationIntervalDays":90,"autoRotationOverlapDays":7,"nextAutoRotationAt":"2024-04-01","rawKey":"wk_secret"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	result, err := client.APIKeys.Create(context.Background(), "org1", "management", CreateApiKeyInput{
		Name: "Key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyId != "k1" {
		t.Errorf("expected KeyId k1, got %q", result.KeyId)
	}
	if result.Type != "management" {
		t.Errorf("expected Type management, got %q", result.Type)
	}
	if result.ScopeId == nil || *result.ScopeId != "org1" {
		t.Errorf("expected ScopeId org1, got %v", result.ScopeId)
	}
	if result.ExpiresAt == nil || *result.ExpiresAt != "2025-01-01" {
		t.Errorf("expected ExpiresAt 2025-01-01, got %v", result.ExpiresAt)
	}
	if result.Attributes == nil || (*result.Attributes)["env"] != "prod" {
		t.Errorf("expected Attributes env=prod, got %v", result.Attributes)
	}
	if result.AutoRotationEnabled == nil || *result.AutoRotationEnabled != true {
		t.Errorf("expected AutoRotationEnabled true")
	}
	if result.RawKey != "wk_secret" {
		t.Errorf("expected RawKey wk_secret, got %q", result.RawKey)
	}
}

// ─── Content-Type Header ──────────────────────────────────────────────────────

func TestRolesResource_Create_ContentTypeHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"r1","name":"Role"}`))
	}))
	defer server.Close()

	client := NewClient(ClientConfig{APIURL: server.URL, APIKey: "test-key"})
	_, err := client.Roles.Create(context.Background(), "scope1", CreateRoleInput{
		ID:            "r1",
		ScopeId:       "scope1",
		Name:          "Role",
		PermissionIds: []string{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ─── float64Ptr helper (not in helpers.go) ─────────────────────────────────────

func float64Ptr(v float64) *float64 { return &v }
