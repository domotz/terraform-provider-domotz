package client

import (
	"fmt"
)

// GetTag retrieves details of a specific tag by listing all tags and filtering
// Note: The API doesn't have a direct GET endpoint for a single tag
func (c *Client) GetTag(tagID int32) (*Tag, error) {
	tags, err := c.ListTags()
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		if tag.ID == tagID {
			return &tag, nil
		}
	}
	return nil, fmt.Errorf("tag with ID %d not found", tagID)
}

// ListTags retrieves all custom tags with pagination
func (c *Client) ListTags() ([]Tag, error) {
	var allTags []Tag
	page := 1
	for {
		path := fmt.Sprintf("/custom-tag?page_size=%d&page_number=%d", defaultPageSize, page)
		var response TagsResponse
		if err := c.doRequest("GET", path, nil, &response); err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}
		allTags = append(allTags, response.Tags...)
		if len(response.Tags) < defaultPageSize {
			break
		}
		page++
	}
	return allTags, nil
}

// CreateTag creates a new custom tag
// Note: API returns 204 No Content, so we need to list tags to find the created one
func (c *Client) CreateTag(req CreateTagRequest) (*Tag, error) {
	path := "/custom-tag"
	if err := c.doRequestNoContent("POST", path, req); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	// API returns 204, so we need to find the created tag by name
	tags, err := c.ListTags()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created tag: %w", err)
	}

	for _, tag := range tags {
		if tag.Name == req.Name {
			return &tag, nil
		}
	}
	return nil, fmt.Errorf("created tag not found")
}

// UpdateTag updates an existing tag
// Note: API returns 204 No Content
func (c *Client) UpdateTag(tagID int32, req UpdateTagRequest) (*Tag, error) {
	path := fmt.Sprintf("/custom-tag/%d", tagID)
	if err := c.doRequestNoContent("PUT", path, req); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	// Retrieve the updated tag
	return c.GetTag(tagID)
}

// DeleteTag deletes a tag
func (c *Client) DeleteTag(tagID int32) error {
	path := fmt.Sprintf("/custom-tag/%d", tagID)
	if err := c.doRequestNoContent("DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

// BindTagToDevice associates a custom tag with a device
func (c *Client) BindTagToDevice(agentID, deviceID, tagID int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d/custom-tag/%d/binding", agentID, deviceID, tagID)
	if err := c.doRequestNoContent("POST", path, nil); err != nil {
		return fmt.Errorf("failed to bind tag to device: %w", err)
	}
	return nil
}

// UnbindTagFromDevice removes a custom tag association from a device
func (c *Client) UnbindTagFromDevice(agentID, deviceID, tagID int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d/custom-tag/%d/binding", agentID, deviceID, tagID)
	if err := c.doRequestNoContent("DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to unbind tag from device: %w", err)
	}
	return nil
}

// ListDeviceTags retrieves all custom tags associated with a device with pagination
func (c *Client) ListDeviceTags(agentID, deviceID int32) ([]Tag, error) {
	var allTags []Tag
	page := 1
	for {
		path := fmt.Sprintf("/agent/%d/device/%d/custom-tag/binding?page_size=%d&page_number=%d",
			agentID, deviceID, defaultPageSize, page)
		var tags []Tag
		if err := c.doRequest("GET", path, nil, &tags); err != nil {
			return nil, fmt.Errorf("failed to list device tags: %w", err)
		}
		allTags = append(allTags, tags...)
		if len(tags) < defaultPageSize {
			break
		}
		page++
	}
	return allTags, nil
}
