package wardenauth

import (
	"context"
	"fmt"
)

func (c *Client) Paginate(ctx context.Context, fn func(ctx context.Context, nextToken *string) (*string, error)) error {
	var nextToken *string

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("pagination cancelled: %w", ctx.Err())
		default:
		}

		result, err := fn(ctx, nextToken)
		if err != nil {
			return err
		}

		if result == nil || *result == "" {
			return nil
		}

		nextToken = result
	}
}
