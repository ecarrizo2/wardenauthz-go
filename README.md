# Go SDK for WardenAuth (Barksoft Access Control)

Go client library for the [WardenAuth](https://wardenauthz.com) managed authorization API.

## Installation

```bash
go get github.com/ecarrizo2/wardenauthz-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/ecarrizo2/wardenauthz-go"
)

func main() {
    client := rbac.NewClient(rbac.ClientConfig{
        APIURL: "https://api.wardenauthz.com",
        APIKey: "sk_your_api_key_here",
    })
    ctx := context.Background()

    // Check access
    result, err := client.Access.HasAccess(ctx, rbac.AccessCheckInput{
        SubjectID: "user-123",
        ScopeID:   "scope-456",
        Resource:  "documents",
        Action:    "read",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Allowed: %v\n", result.Allowed)

    // List scopes
    scopes, err := client.Scopes.List(ctx, nil)
    if err != nil {
        panic(err)
    }
    for _, s := range scopes.Items {
        fmt.Printf("Scope: %s (%s)\n", s.Name, s.Type)
    }
}
```

## Resources

The client exposes 11 resource groups:

| Resource | Description |
|----------|-------------|
| `client.Access` | Access evaluation (check, simulate, list permissions/roles) and permission **receipts** (issue/verify) |
| `client.Scopes` | Scope management (CRUD, manifest apply, clone) |
| `client.Permissions` | Permission management (CRUD, bulk, CSV import) |
| `client.Roles` | Role management (CRUD, bulk, clone) |
| `client.AccessPolicies` | Access policy assignment (CRUD, CSV import) |
| `client.APIKeys` | API key management (CRUD, rotate, reveal, auto-rotation) |
| `client.Webhooks` | Webhook endpoint management (CRUD, rotate secret, test, deliveries) |
| `client.Audit` | Audit log querying, export (CSV/JSON), and signature verification |
| `client.SessionTokens` | Session token minting (mint, intent-mint, verify-intent-call, revoke) |
| `client.SodConstraints` | Separation of duty constraints (CRUD) |
| `client.TeamMembers` | Team member management (add, list, remove) |
| `client.ResourceTypes` | Resource type management (CRUD) |
| `client.Tuples` | Relationship tuple management (write, list, listByResource) |

## Error Handling

All methods return errors. API errors are typed as `*rbac.APIError`:

```go
result, err := client.Scopes.Create(ctx, input)
if err != nil {
    if apiErr, ok := err.(*rbac.APIError); ok {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, string(apiErr.Body))
    }
    return
}
```

## Automatic Retries

All HTTP methods automatically retry on transient errors (429, 500, 502, 503, 504) with exponential backoff. Configure retry behavior via `ClientConfig`:

```go
// Disable retries
client := rbac.NewClient(rbac.ClientConfig{
    APIURL:     "https://api.wardenauthz.com",
    APIKey:     "sk_your_api_key",
    MaxRetries: 0, // disables retry
})

// Custom retry settings
client := rbac.NewClient(rbac.ClientConfig{
    APIURL:     "https://api.wardenauthz.com",
    APIKey:     "sk_your_api_key",
    MaxRetries: 5,                      // up to 5 retries (default: 3)
    RetryDelay: 1 * time.Second,        // initial backoff (default: 500ms)
    Timeout:    10 * time.Second,       // per-request timeout (default: 30s)
})
```

**Default retry behavior:**
- Retries on status codes: `429`, `500`, `502`, `503`, `504`
- Max 3 retry attempts
- Exponential backoff starting at 500ms (500ms, 1s, 2s)
- Capped at 30s max delay
- Non-retryable errors (4xx except 429) are thrown immediately
- Retries respect `context.Context` cancellation

## Helper Functions

The package provides convenience pointer functions for optional fields:

```go
import "github.com/ecarrizo2/wardenauthz-go"

input := rbac.CreateAccessPolicyInput{
    SubjectID:  "user-123",
    RoleIDs:    []string{"editor"},
    ExpiresAt:  rbac.StringPtr("2026-12-31T23:59:59Z"),
}

opts := rbac.ListOptions{
    Limit: rbac.IntPtr(100),
}
```

Available helpers:

| Function              | Returns     | Use for                |
| --------------------- | ----------- | ---------------------- |
| `rbac.StringPtr(s)`   | `*string`   | Optional string fields |
| `rbac.IntPtr(n)`      | `*int`      | Optional integer fields |
| `rbac.BoolPtr(b)`     | `*bool`     | Optional boolean fields |

## Pagination

The client provides a `Paginate` helper for cursor-based list endpoints. It handles the iteration automatically:

```go
// Iterate over all permissions in a scope
err := client.Paginate(ctx, func(ctx context.Context, nextToken *string) (*string, error) {
    opts := &struct{ Limit, NextToken string }{
        Limit: "50",
    }
    if nextToken != nil {
        opts.NextToken = *nextToken
    }

    result, err := client.Permissions.List(ctx, "workspace-abc", opts)
    if err != nil {
        return nil, err
    }

    for _, p := range result.Items {
        fmt.Printf("Permission: %s (%s:%s)\n", p.ID, p.Resource, p.Action)
    }

    return result.NextToken, nil
})
if err != nil {
    log.Fatal(err)
}
```

`Paginate` respects `context.Context` — a cancelled context stops iteration immediately.

## Audit

Query, export, and verify audit logs:

```go
// List audit events
result, err := client.Audit.List(ctx, &struct{
    ScopeID string
    Limit   string
}{ScopeID: "org-abc", Limit: "50"})

// Export as CSV
csv, err := client.Audit.Export(ctx, &struct{
    ScopeID string
    Format  rbac.AuditExportFormat
}{ScopeID: "org-abc", Format: rbac.AuditExportCSV})

// Verify HMAC signatures (SOC2 compliance)
verifyResult, err := client.Audit.Verify(ctx, rbac.AuditVerifyInput{
    ScopeID:    "org-abc",
    StartDate:  "2026-01-01",
    EndDate:    "2026-06-01",
    MaxRecords: 5000,
})
// verifyResult.Matched, verifyResult.MismatchCount, ...
```

## Context Support

All methods accept a `context.Context` for cancellation and deadline support:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
result, err := client.Access.HasAccess(ctx, input)
```

## Connection Pooling

The client uses a tuned `http.Transport` with keep-alive connection reuse
(`MaxIdleConnsPerHost: 64`, HTTP/2 enabled) so concurrent calls to the API reuse warm
TLS connections instead of paying a handshake per request. Reuse a single `*Client`
across goroutines to benefit from the pool.

## License

Proprietary — see the main repository for details.
