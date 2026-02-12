package client

import (
	"context"
	"fmt"
)

// GetDevice retrieves details of a specific device
func (c *Client) GetDevice(ctx context.Context, agentID, deviceID int32) (*Device, error) {
	path := fmt.Sprintf("/agent/%d/device/%d", agentID, deviceID)
	var device Device
	if err := c.doRequest(ctx, "GET", path, nil, &device); err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return &device, nil
}

// ListDevices retrieves all devices for a specific agent
func (c *Client) ListDevices(ctx context.Context, agentID int32) ([]Device, error) {
	path := fmt.Sprintf("/agent/%d/device", agentID)
	var devices []Device
	if err := c.doRequest(ctx, "GET", path, nil, &devices); err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	return devices, nil
}

// CreateDevice creates a new external IP device (external host)
func (c *Client) CreateDevice(ctx context.Context, agentID int32, req CreateDeviceRequest) (*Device, error) {
	path := fmt.Sprintf("/agent/%d/device/external-host", agentID)
	var device Device
	if err := c.doRequest(ctx, "POST", path, req, &device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}
	return &device, nil
}

// UpdateDeviceImportance updates the importance level of a device
func (c *Client) UpdateDeviceImportance(ctx context.Context, agentID, deviceID int32, importance string) error {
	path := fmt.Sprintf("/agent/%d/device/%d/importance", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "PUT", path, importance); err != nil {
		return fmt.Errorf("failed to update device importance: %w", err)
	}
	return nil
}

// UpdateDeviceUserDataName updates the user_data name field of a device
func (c *Client) UpdateDeviceUserDataName(ctx context.Context, agentID, deviceID int32, name string) error {
	path := fmt.Sprintf("/agent/%d/device/%d/user_data/name", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "PUT", path, name); err != nil {
		return fmt.Errorf("failed to update device user_data name: %w", err)
	}
	return nil
}

// UpdateDeviceUserDataModel updates the user_data model field of a device
func (c *Client) UpdateDeviceUserDataModel(ctx context.Context, agentID, deviceID int32, model string) error {
	path := fmt.Sprintf("/agent/%d/device/%d/user_data/model", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "PUT", path, model); err != nil {
		return fmt.Errorf("failed to update device user_data model: %w", err)
	}
	return nil
}

// UpdateDeviceUserDataVendor updates the user_data vendor field of a device
func (c *Client) UpdateDeviceUserDataVendor(ctx context.Context, agentID, deviceID int32, vendor string) error {
	path := fmt.Sprintf("/agent/%d/device/%d/user_data/vendor", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "PUT", path, vendor); err != nil {
		return fmt.Errorf("failed to update device user_data vendor: %w", err)
	}
	return nil
}

// UpdateDeviceUserDataType updates the user_data type field of a device
func (c *Client) UpdateDeviceUserDataType(ctx context.Context, agentID, deviceID, deviceType int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d/user_data/type", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "PUT", path, deviceType); err != nil {
		return fmt.Errorf("failed to update device user_data type: %w", err)
	}
	return nil
}

// UpdateDevice updates an existing device by calling individual field update endpoints
// This is a convenience method that calls the appropriate field-specific endpoints
func (c *Client) UpdateDevice(ctx context.Context, agentID, deviceID int32, req UpdateDeviceRequest) (*Device, error) {
	// Update importance if provided
	if req.Importance != nil && *req.Importance != "" {
		if err := c.UpdateDeviceImportance(ctx, agentID, deviceID, *req.Importance); err != nil {
			return nil, err
		}
	}

	// Update user_data fields if provided
	if req.UserData != nil {
		if req.UserData.Name != "" {
			if err := c.UpdateDeviceUserDataName(ctx, agentID, deviceID, req.UserData.Name); err != nil {
				return nil, err
			}
		}
		if req.UserData.Model != "" {
			if err := c.UpdateDeviceUserDataModel(ctx, agentID, deviceID, req.UserData.Model); err != nil {
				return nil, err
			}
		}
		if req.UserData.Vendor != "" {
			if err := c.UpdateDeviceUserDataVendor(ctx, agentID, deviceID, req.UserData.Vendor); err != nil {
				return nil, err
			}
		}
		if req.UserData.Type != 0 {
			if err := c.UpdateDeviceUserDataType(ctx, agentID, deviceID, req.UserData.Type); err != nil {
				return nil, err
			}
		}
	}

	// Return the updated device
	return c.GetDevice(ctx, agentID, deviceID)
}

// DeleteDevice deletes a device
func (c *Client) DeleteDevice(ctx context.Context, agentID, deviceID int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d", agentID, deviceID)
	if err := c.doRequestNoContent(ctx, "DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}
