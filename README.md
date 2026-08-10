# WardenAuthz Go SDK

[![CI](https://github.com/ecarrizo2/wardenauthz-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ecarrizo2/wardenauthz-go/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/ecarrizo2/wardenauthz-go.svg)](https://pkg.go.dev/github.com/ecarrizo2/wardenauthz-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go client for the [WardenAuthz](https://wardenauthz.com) managed authorization API.

## Installation

```bash
go get github.com/ecarrizo2/wardenauthz-go
```

Requires Go 1.26+.

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    wardenauth "github.com/ecarrizo2/wardenauthz-go"
)

func main() {
    client := wardenauth.NewClient(wardenauth.ClientConfig{
        APIURL: "https://api.wardenauthz.com",
        APIKey: "sk_your_api_key_here",
    })
    ctx := context.Background()

    // Check access
    result, err := client.Access.HasAccess(ctx, wardenauth.AccessCheckInput{
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

All methods accept `context.Context` as the first parameter for cancellation and deadline support.

### Scopes (`client.Scopes`)

| Method | Description |
|--------|-------------|
| `Create(ctx, input)` | Create a new scope |
| `List(ctx, opts)` | List scopes with optional pagination and type filter |
| `GetByID(ctx, scopeID)` | Get a scope by ID |
| `Update(ctx, scopeID, input)` | Update a scope |
| `Delete(ctx, scopeID)` | Delete a scope |
| `ApplyManifest(ctx, scopeID, input)` | Apply an authorization manifest to a scope |
| `Export(ctx, scopeID)` | Export a scope's authorization configuration |

### Permissions (`client.Permissions`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create a permission in a scope |
| `BulkCreate(ctx, scopeID, inputs)` | Create multiple permissions at once |
| `List(ctx, scopeID, opts)` | List permissions in a scope (paginated) |
| `GetByID(ctx, scopeID, permissionID)` | Get a permission by ID |
| `Update(ctx, scopeID, permissionID, input)` | Update a permission |
| `Delete(ctx, scopeID, permissionID)` | Delete a permission |
| `BulkDelete(ctx, scopeID, ids)` | Delete multiple permissions at once |
| `ImportCSV(ctx, scopeID, csvContent)` | Import permissions from CSV content |

### Roles (`client.Roles`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create a role in a scope |
| `BulkCreate(ctx, scopeID, inputs)` | Create multiple roles at once |
| `List(ctx, scopeID, opts)` | List roles in a scope (paginated) |
| `GetByID(ctx, scopeID, roleID)` | Get a role by ID |
| `Update(ctx, scopeID, roleID, input)` | Update a role |
| `Delete(ctx, scopeID, roleID)` | Delete a role |
| `BulkDelete(ctx, scopeID, ids)` | Delete multiple roles at once |
| `Clone(ctx, targetScopeID, templateRoleID)` | Clone a role from another scope |

### Access Policies (`client.AccessPolicies`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create an access policy in a scope |
| `ListByScope(ctx, scopeID)` | List all access policies in a scope |
| `ListBySubject(ctx, scopeID, subjectID)` | List access policies for a subject |
| `GetByID(ctx, scopeID, policyID)` | Get an access policy by ID |
| `Update(ctx, scopeID, policyID, input)` | Update an access policy |
| `Delete(ctx, scopeID, policyID)` | Delete an access policy |
| `ImportCSV(ctx, scopeID, csvContent)` | Import access policies from CSV |

### API Keys (`client.APIKeys`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, keyType, input)` | Create an API key (`management` or `application`) |
| `List(ctx, scopeID, keyType)` | List API keys of a given type |
| `GetByID(ctx, scopeID, keyType, keyID)` | Get an API key by ID |
| `Delete(ctx, scopeID, keyType, keyID)` | Delete an API key |
| `Rotate(ctx, scopeID, keyType, keyID)` | Rotate an API key |
| `RevealRotation(ctx, scopeID, keyType, keyID, rotationRef)` | Reveal a rotated API key value |
| `UpdateAutoRotation(ctx, scopeID, keyType, keyID, input)` | Configure automatic key rotation |

### Webhooks (`client.Webhooks`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create a webhook endpoint |
| `List(ctx, scopeID, opts)` | List webhook endpoints (paginated) |
| `GetByID(ctx, scopeID, endpointID)` | Get a webhook endpoint by ID |
| `Update(ctx, scopeID, endpointID, input)` | Update a webhook endpoint |
| `Delete(ctx, scopeID, endpointID)` | Delete a webhook endpoint |
| `RotateSecret(ctx, scopeID, endpointID)` | Rotate the webhook signing secret |
| `Test(ctx, scopeID, endpointID)` | Send a test event to the endpoint |
| `RetryDelivery(ctx, scopeID, endpointID, deliveryID)` | Retry a failed delivery |
| `ListDeliveries(ctx, scopeID, endpointID, opts)` | List delivery attempts (paginated) |

### Access (`client.Access`)

| Method | Description |
|--------|-------------|
| `HasAccess(ctx, input)` | Check if a subject has access to a resource+action |
| `HasAccessBulk(ctx, inputs)` | Check multiple access decisions at once |
| `ListPermissions(ctx, input)` | List effective permissions for a subject |
| `ListRoles(ctx, input)` | List effective roles for a subject |
| `IssueReceipt(ctx, input)` | Issue a cryptographically signed access receipt |
| `VerifyReceipt(ctx, input)` | Verify a receipt's validity |
| `Simulate(ctx, input)` | Simulate access with hypothetical permissions/roles |

### Audit (`client.Audit`)

| Method | Description |
|--------|-------------|
| `List(ctx, opts)` | Query audit log events (paginated) |
| `Export(ctx, opts)` | Export audit logs as CSV or JSON |
| `Verify(ctx, input)` | Verify audit log integrity (HMAC) |

### Session Tokens (`client.SessionTokens`)

| Method | Description |
|--------|-------------|
| `Mint(ctx, input)` | Mint a session token with specific permissions |
| `MintWithIntent(ctx, input)` | Mint a session token with an agent intent |
| `VerifyIntentCall(ctx, input)` | Verify an intent-scoped session token call |
| `Revoke(ctx, jti)` | Revoke a session token by JTI |

### Separation of Duty (`client.SodConstraints`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create an SoD constraint |
| `List(ctx, scopeID)` | List SoD constraints in a scope |
| `Delete(ctx, scopeID, constraintID)` | Delete an SoD constraint |

### Team Members (`client.TeamMembers`)

| Method | Description |
|--------|-------------|
| `Add(ctx, scopeID, input)` | Add a team member to a scope |
| `List(ctx, scopeID)` | List team members in a scope |
| `Remove(ctx, scopeID, subjectID)` | Remove a team member |

### Resource Types (`client.ResourceTypes`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create a resource type |
| `List(ctx, scopeID, opts)` | List resource types (paginated) |
| `Delete(ctx, scopeID, typeID)` | Delete a resource type |

### Tuples (`client.Tuples`)

| Method | Description |
|--------|-------------|
| `Write(ctx, scopeID, input)` | Write relationship tuples (create/delete) |
| `List(ctx, scopeID, subjectID, resourceType, resourceID)` | List tuples by subject |
| `ListByResource(ctx, scopeID, resourceType, resourceID)` | List tuples by resource |

### Integrations (`client.Integrations`)

| Method | Description |
|--------|-------------|
| `Create(ctx, scopeID, input)` | Create an integration |
| `List(ctx, scopeID)` | List integrations in a scope |
| `GetByID(ctx, scopeID, id)` | Get an integration by ID |
| `Update(ctx, scopeID, id, input)` | Update an integration |
| `Delete(ctx, scopeID, id)` | Delete an integration |

### MCP (`client.MCP`)

| Method | Description |
|--------|-------------|
| `ListAssignments(ctx, scopeID, subjectID)` | List MCP server assignments for a user |
| `SetAssignment(ctx, scopeID, subjectID, serverKey, maxTier, constraints)` | Grant/update MCP server access |
| `DeleteAssignment(ctx, scopeID, subjectID, serverKey)` | Revoke MCP server access |

### Consent (`client.Consent`)

| Method | Description |
|--------|-------------|
| `GetServers(ctx, scopeID)` | Get available MCP servers for consent |
| `GetContext(ctx, serverKey, scopeID)` | Get consent context for a server |
| `GetPortalContext(ctx, scopeID)` | Get portal consent context |
| `Grant(ctx, input)` | Grant consent for MCP server access |
| `Deny(ctx, authRequestID)` | Deny a consent request |
| `GetRequestInfo(ctx, reqID)` | Get details of a consent request |
| `ListGrants(ctx)` | List all active consent grants |
| `RevokeGrant(ctx, grantID)` | Revoke a consent grant |
| `ListApprovals(ctx)` | List pending MCP approvals |
| `ListApprovalHistory(ctx)` | List approval history |
| `ApproveRequest(ctx, id, reason)` | Approve a pending approval request |
| `DenyRequest(ctx, id, reason)` | Deny a pending approval request |
| `GetPushPublicKey(ctx)` | Get the push notification public key |
| `SubscribePush(ctx, endpoint, p256dh, auth)` | Subscribe to push notifications |
| `UnsubscribePush(ctx, endpoint)` | Unsubscribe from push notifications |
| `GetVelocityConfig(ctx)` | Get rate-limiting velocity config |
| `UpdateVelocityConfig(ctx, input)` | Update rate-limiting velocity config |

### Agent (`client.Agent`)

| Method | Description |
|--------|-------------|
| `Identify(ctx, scopeID, input)` | Identify an AI agent with delegated permissions |
| `Check(ctx, scopeID, input)` | Check if an agent's action is authorized |

## Error Handling

All methods return errors. API errors are typed as `*wardenauth.APIError`:

```go
result, err := client.Scopes.Create(ctx, input)
if err != nil {
    if apiErr, ok := err.(*wardenauth.APIError); ok {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, string(apiErr.Body))
    }
    return
}
```

The `APIError` type provides:

```go
type APIError struct {
    StatusCode int             // HTTP status code
    Body       json.RawMessage // Raw response body
    Message    string          // HTTP status text
}
```

## Retry Behavior

All HTTP methods automatically retry on transient errors with exponential backoff.

**Retryable status codes:** `429` (Too Many Requests), `500`, `502`, `503`, `504`, and any `5xx`.

**Non-retryable status codes:** All `4xx` errors except `429` are returned immediately without retrying.

Configuration:

```go
client := wardenauth.NewClient(wardenauth.ClientConfig{
    APIURL:     "https://api.wardenauthz.com",
    APIKey:     "sk_your_api_key",
    MaxRetries: 3,                      // max retry attempts (default: 3)
    RetryDelay: 500 * time.Millisecond, // initial backoff (default: 500ms)
    Timeout:    30 * time.Second,       // per-request timeout (default: 30s)
})
```

**Default behavior:**
- Up to 3 retry attempts
- Exponential backoff: 500ms, 1s, 2s, 4s, ...
- Capped at 30s maximum delay
- Respects `context.Context` cancellation during backoff

## Advanced Configuration

### Connection Pooling

The client uses a tuned `http.Transport` with HTTP/2 enabled and keep-alive connection reuse (`MaxIdleConnsPerHost: 64`). Reuse a single `*Client` across goroutines:

```go
var client = wardenauth.NewClient(wardenauth.ClientConfig{...})

// Safe for concurrent use
go func() { client.Access.HasAccess(ctx, input) }()
go func() { client.Scopes.List(ctx, nil) }()
```

### Custom Timeouts

```go
client := wardenauth.NewClient(wardenauth.ClientConfig{
    APIURL:  "https://api.wardenauthz.com",
    APIKey:  "sk_key",
    Timeout: 5 * time.Second, // per-request timeout
})

// Per-call deadlines via context
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
result, err := client.Access.HasAccess(ctx, input)
```

### Disabling Retries

```go
client := wardenauth.NewClient(wardenauth.ClientConfig{
    APIURL:     "https://api.wardenauthz.com",
    APIKey:     "sk_key",
    MaxRetries: 0, // disables retry
})
```

## Pagination

Use the `Paginate` helper for cursor-based list endpoints:

```go
err := client.Paginate(ctx, func(ctx context.Context, nextToken *string) (*string, error) {
    result, err := client.Permissions.List(ctx, "workspace-abc", &struct {
        Limit, NextToken string
    }{
        Limit: "50",
        NextToken: func() string {
            if nextToken != nil { return *nextToken }
            return ""
        }(),
    })
    if err != nil {
        return nil, err
    }

    for _, p := range result.Items {
        fmt.Printf("Permission: %s (%s:%s)\n", p.ID, p.Resource, p.Action)
    }

    if result.NextToken == "" {
        return nil, nil
    }
    return &result.NextToken, nil
})
```

`Paginate` respects `context.Context` — a cancelled context stops iteration immediately.

## Helper Functions

The package provides convenience pointer functions for optional fields:

```go
input := wardenauth.CreateAccessPolicyInput{
    SubjectID: "user-123",
    ScopeID:   "scope-456",
    RoleIds:   &[]string{"editor"},
    ExpiresAt: wardenauth.StringPtr("2026-12-31T23:59:59Z"),
}
```

| Function | Returns | Use for |
|----------|---------|---------|
| `wardenauth.StringPtr(s)` | `*string` | Optional string fields |
| `wardenauth.IntPtr(n)` | `*int` | Optional integer fields |
| `wardenauth.BoolPtr(b)` | `*bool` | Optional boolean fields |

## API Reference

Full API documentation is available at the [WardenAuthz documentation portal](https://wardenauthz.com/docs).

## License

Proprietary — see the main repository for details.
