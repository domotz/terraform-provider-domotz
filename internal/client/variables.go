package client

import "fmt"

// GetVariable retrieves details of a specific variable
func (c *Client) GetVariable(agentID, deviceID, variableID int32) (*Variable, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/variable/%d", agentID, deviceID, variableID)
	var variable Variable
	if err := c.doRequest("GET", path, nil, &variable); err != nil {
		return nil, fmt.Errorf("failed to get variable: %w", err)
	}
	return &variable, nil
}

// ListVariables retrieves all variables for a device
func (c *Client) ListVariables(agentID, deviceID int32) ([]Variable, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/variable", agentID, deviceID)
	var variables []Variable
	if err := c.doRequest("GET", path, nil, &variables); err != nil {
		return nil, fmt.Errorf("failed to list variables: %w", err)
	}
	return variables, nil
}
