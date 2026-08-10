package wardenauth

import (
	"context"
	"fmt"
)

type TuplesResource struct{ client *Client }

func (r *TuplesResource) Write(ctx context.Context, scopeID string, input TupleWriteInput) (*TupleWriteResult, error) {
	var result TupleWriteResult
	if err := r.client.post(ctx, fmt.Sprintf("/v1/scope/%s/tuples", scopeID), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *TuplesResource) List(
	ctx context.Context,
	scopeID string,
	subjectID string,
	resourceType *string,
	resourceID *string,
) (*TupleListResult, error) {
	params := map[string]string{"subjectId": subjectID}
	if resourceType != nil {
		params["resourceType"] = *resourceType
	}
	if resourceID != nil {
		params["resourceId"] = *resourceID
	}
	path := buildQueryPath(fmt.Sprintf("/v1/scope/%s/tuples", scopeID), params)

	var result TupleListResult
	if err := r.client.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *TuplesResource) ListByResource(
	ctx context.Context,
	scopeID string,
	resourceType string,
	resourceID string,
) (*TupleListResult, error) {
	params := map[string]string{"resourceType": resourceType, "resourceId": resourceID}
	path := buildQueryPath(fmt.Sprintf("/v1/scope/%s/tuples-by-resource", scopeID), params)

	var result TupleListResult
	if err := r.client.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
