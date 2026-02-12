package client

import (
	"fmt"
)

// GetSNMPSensor retrieves details of a specific SNMP sensor (Domotz Eye)
// Note: The API doesn't have a direct GET for a single sensor, so we list and filter
func (c *Client) GetSNMPSensor(agentID, deviceID, sensorID int32) (*SNMPSensor, error) {
	sensors, err := c.ListSNMPSensors(agentID, deviceID)
	if err != nil {
		return nil, err
	}
	for _, sensor := range sensors {
		if sensor.ID == sensorID {
			return &sensor, nil
		}
	}
	return nil, fmt.Errorf("SNMP sensor with ID %d not found", sensorID)
}

// ListSNMPSensors retrieves all SNMP sensors (Domotz Eyes) for a device
func (c *Client) ListSNMPSensors(agentID, deviceID int32) ([]SNMPSensor, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/snmp", agentID, deviceID)
	var sensors []SNMPSensor
	if err := c.doRequest("GET", path, nil, &sensors); err != nil {
		return nil, fmt.Errorf("failed to list SNMP sensors: %w", err)
	}
	return sensors, nil
}

// CreateSNMPSensor creates a new SNMP sensor (Domotz Eye)
// Note: API returns 201 with empty body, so we list and find by OID
func (c *Client) CreateSNMPSensor(agentID, deviceID int32, req CreateSNMPSensorRequest) (*SNMPSensor, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/snmp", agentID, deviceID)
	if err := c.doRequestNoContent("POST", path, req); err != nil {
		return nil, fmt.Errorf("failed to create SNMP sensor: %w", err)
	}
	// API returns empty body, find created sensor by OID
	sensors, err := c.ListSNMPSensors(agentID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find created sensor: %w", err)
	}
	for _, s := range sensors {
		if s.OID == req.OID {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("created sensor not found")
}

// DeleteSNMPSensor deletes an SNMP sensor (Domotz Eye)
func (c *Client) DeleteSNMPSensor(agentID, deviceID, sensorID int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/snmp/%d", agentID, deviceID, sensorID)
	if err := c.doRequestNoContent("DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to delete SNMP sensor: %w", err)
	}
	return nil
}

// GetTCPSensor retrieves details of a specific TCP sensor (Domotz Eye)
// Note: The API doesn't have a direct GET for a single sensor, so we list and filter
func (c *Client) GetTCPSensor(agentID, deviceID, sensorID int32) (*TCPSensor, error) {
	sensors, err := c.ListTCPSensors(agentID, deviceID)
	if err != nil {
		return nil, err
	}
	for _, sensor := range sensors {
		if sensor.ID == sensorID {
			return &sensor, nil
		}
	}
	return nil, fmt.Errorf("TCP sensor with ID %d not found", sensorID)
}

// ListTCPSensors retrieves all TCP sensors (Domotz Eyes) for a device
func (c *Client) ListTCPSensors(agentID, deviceID int32) ([]TCPSensor, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/tcp", agentID, deviceID)
	var sensors []TCPSensor
	if err := c.doRequest("GET", path, nil, &sensors); err != nil {
		return nil, fmt.Errorf("failed to list TCP sensors: %w", err)
	}
	return sensors, nil
}

// CreateTCPSensor creates a new TCP sensor (Domotz Eye)
// Note: API returns 201 with empty body, so we list and find by port
func (c *Client) CreateTCPSensor(agentID, deviceID int32, req CreateTCPSensorRequest) (*TCPSensor, error) {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/tcp", agentID, deviceID)
	if err := c.doRequestNoContent("POST", path, req); err != nil {
		return nil, fmt.Errorf("failed to create TCP sensor: %w", err)
	}
	// API returns empty body, find created sensor by port
	sensors, err := c.ListTCPSensors(agentID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find created sensor: %w", err)
	}
	for _, s := range sensors {
		if s.Port == req.Port {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("created sensor not found")
}

// DeleteTCPSensor deletes a TCP sensor (Domotz Eye)
func (c *Client) DeleteTCPSensor(agentID, deviceID, sensorID int32) error {
	path := fmt.Sprintf("/agent/%d/device/%d/eye/tcp/%d", agentID, deviceID, sensorID)
	if err := c.doRequestNoContent("DELETE", path, nil); err != nil {
		return fmt.Errorf("failed to delete TCP sensor: %w", err)
	}
	return nil
}
