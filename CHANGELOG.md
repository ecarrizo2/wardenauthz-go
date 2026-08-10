# Changelog

All notable changes to the Go SDK (`github.com/ecarrizo/warden-auth-go`) are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and
this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0] - 2026-06-26

### Added

- `TierPolicyResource` (`client.TierPolicy`) — `Get()`, `Update()` for per-scope agent tool trust-tier policy (`allow` / `approve` / `deny`)
- `AccessResource` — `IssueReceipt()`, `VerifyReceipt()` for HMAC-signed permission receipts (non-repudiation)

### Changed

- HTTP transport — the client now uses a tuned `http.Transport` (`MaxIdleConnsPerHost: 64`, HTTP/2 enabled) so concurrent calls reuse warm keep-alive connections instead of paying a TLS handshake per request

## [v0.1.0] - Initial Release

### Added

- `Client` with resource registry: `Scopes`, `Permissions`, `Roles`, `AccessPolicies`, `APIKeys`, `Webhooks`, `Access`, `Audit`, `Billing`, `Dashboard`, `SessionTokens`, `SodConstraints`, `SSO`, `TeamMembers`, `ResourceTypes`, `Organization`, `Identities`, `Tuples`, `Integrations`
- Context-aware methods, typed errors (`APIError`), automatic retries, and pagination helpers
