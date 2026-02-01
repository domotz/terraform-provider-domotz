package client

import "fmt"

// GetAgent retrieves details of a specific agent
func (c *Client) GetAgent(agentID int32) (*Agent, error) {
	path := fmt.Sprintf("/agent/%d", agentID)
	var agent Agent
	if err := c.doRequest("GET", path, nil, &agent); err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return &agent, nil
}

// ListAgents retrieves a list of all agents
func (c *Client) ListAgents() ([]Agent, error) {
	path := "/agent"
	var agents []Agent
	if err := c.doRequest("GET", path, nil, &agents); err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}
	return agents, nil
}

// DeleteAgent deletes an agent
func (c *Client) DeleteAgent(agentID int32) error {
	path := fmt.Sprintf("/agent/%d", agentID)
	if err := c.doRequestNoContent("DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}
	return nil
}
