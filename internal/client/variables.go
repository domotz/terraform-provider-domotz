package client

import (
	"context"
	"fmt"
)

// GetVariable retrieves details of a specific variable
func (c *Client) GetVariable(ctx context.Context, agentID, deviceID, variableID int32) (*Variable, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/variable/%d", agentID, deviceID, variableID)
	var variable Variable
	if err := c.doRequest(ctx, "GET", path, nil, &variable); err != nil {
		return nil, fmt.Errorf("failed to get variable: %w", err)
	}
	return &variable, nil
}

// ListVariables retrieves all variables for a device with pagination
func (c *Client) ListVariables(ctx context.Context, agentID, deviceID int32) ([]Variable, error) {
	var allVariables []Variable
	page := 1
	for {
		path := fmt.Sprintf("/agent/%d/device/%d/variable?page_size=%d&page_number=%d",
			agentID, deviceID, defaultPageSize, page)
		var variables []Variable
		if err := c.doRequest(ctx, "GET", path, nil, &variables); err != nil {
			return nil, fmt.Errorf("failed to list variables: %w", err)
		}
		allVariables = append(allVariables, variables...)
		if len(variables) < defaultPageSize {
			break
		}
		page++
	}
	return allVariables, nil
}
