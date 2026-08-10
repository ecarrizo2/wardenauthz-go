package wardenauth

import (
	"context"
	"fmt"
)

type AuditResource struct{ client *Client }

func (r *AuditResource) List(ctx context.Context, opts *struct {
	ScopeID      string
	Limit        string
	NextToken    string
	ActorID      string
	ResourceType string
	ResourceID   string
	StartDate    string
	EndDate      string
}) (*PaginatedResult[AuditLogItem], error) {
	params := map[string]string{}
	if opts != nil {
		if opts.ScopeID != "" { params["scopeId"] = opts.ScopeID }
		if opts.Limit != "" { params["limit"] = opts.Limit }
		if opts.NextToken != "" { params["nextToken"] = opts.NextToken }
		if opts.ActorID != "" { params["actorId"] = opts.ActorID }
		if opts.ResourceType != "" { params["resourceType"] = opts.ResourceType }
		if opts.ResourceID != "" { params["resourceId"] = opts.ResourceID }
		if opts.StartDate != "" { params["startDate"] = opts.StartDate }
		if opts.EndDate != "" { params["endDate"] = opts.EndDate }
	}
	var result PaginatedResult[AuditLogItem]
	if err := r.client.get(ctx, buildQueryPath("/v1/audit", params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AuditResource) Export(ctx context.Context, opts *struct {
	ScopeID      string
	ActorID      string
	ResourceType string
	ResourceID   string
	StartDate    string
	EndDate      string
	EventTypes   []string
	MaxRows      string
	Format       AuditExportFormat
}) (string, error) {
	params := map[string]string{}
	if opts != nil {
		if opts.ScopeID != "" { params["scopeId"] = opts.ScopeID }
		if opts.ActorID != "" { params["actorId"] = opts.ActorID }
		if opts.ResourceType != "" { params["resourceType"] = opts.ResourceType }
		if opts.ResourceID != "" { params["resourceId"] = opts.ResourceID }
		if opts.StartDate != "" { params["startDate"] = opts.StartDate }
		if opts.EndDate != "" { params["endDate"] = opts.EndDate }
		if opts.MaxRows != "" { params["maxRows"] = opts.MaxRows }
		if opts.Format != "" { params["format"] = string(opts.Format) }
	}
	query := buildQueryPath("/v1/audit/export", params)
	if opts != nil && len(opts.EventTypes) > 0 {
		for _, et := range opts.EventTypes {
			query += fmt.Sprintf("&eventTypes=%s", et)
		}
	}
	return r.client.getRaw(ctx, query)
}

func (r *AuditResource) Verify(ctx context.Context, input AuditVerifyInput) (*AuditVerifyResult, error) {
	params := map[string]string{"scopeId": input.ScopeId}
	if input.StartDate != nil && *input.StartDate != "" {
		params["startDate"] = *input.StartDate
	}
	if input.EndDate != nil && *input.EndDate != "" {
		params["endDate"] = *input.EndDate
	}
	if input.MaxRecords != nil && *input.MaxRecords > 0 {
		params["maxRecords"] = fmt.Sprintf("%d", int(*input.MaxRecords))
	}

	var result AuditVerifyResult
	if err := r.client.get(ctx, buildQueryPath("/v1/audit/verify", params), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
