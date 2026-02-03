package client

import "time"

// Team represents the team/area associated with an agent
type Team struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// AgentStatus represents the status of an agent
type AgentStatus struct {
	Value      string    `json:"value"`       // ONLINE, OFFLINE
	LastChange time.Time `json:"last_change"` // When status last changed
}

// Agent represents a Domotz collector/agent
type Agent struct {
	ID           int32       `json:"id"`
	DisplayName  string      `json:"display_name"`
	Status       AgentStatus `json:"status"`
	Team         Team        `json:"team"`
	CreationTime time.Time   `json:"creation_time"`
	OnlineAt     time.Time   `json:"online_at,omitempty"`
}

// DeviceUserData represents custom metadata for a device
type DeviceUserData struct {
	Name   string `json:"name,omitempty"`
	Model  string `json:"model,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Type   int32  `json:"type,omitempty"`
}

// Device represents a monitored device
type Device struct {
	ID                   int32          `json:"id"`
	AgentID              int32          `json:"agent_id"`
	DisplayName          string         `json:"display_name"`
	Protocol             string         `json:"protocol"` // IP, DUMMY, etc.
	IPAddresses          []string       `json:"ip_addresses,omitempty"`
	Vendor               string         `json:"vendor,omitempty"` // Auto-discovered vendor (e.g., "Ubiquiti Inc")
	Model                string         `json:"model,omitempty"`  // Auto-discovered model
	UserData             DeviceUserData `json:"user_data"`        // User-editable metadata
	AuthenticationStatus string         `json:"authentication_status,omitempty"`
	Importance           string         `json:"importance,omitempty"` // VITAL, FLOATING
	HWAddress            string         `json:"hw_address,omitempty"`
	Zone                 string         `json:"zone,omitempty"`
	FirstSeenAt          time.Time      `json:"first_seen_at,omitempty"`
	LastStatusChange     time.Time      `json:"last_status_change,omitempty"`
}

// CreateDeviceRequest represents the request to create a new device
type CreateDeviceRequest struct {
	DisplayName string         `json:"display_name"`
	IPAddresses []string       `json:"ip_addresses"`
	UserData    DeviceUserData `json:"user_data,omitempty"`
	Importance  string         `json:"importance,omitempty"`
}

// UpdateDeviceRequest represents the request to update a device
type UpdateDeviceRequest struct {
	DisplayName *string         `json:"display_name,omitempty"`
	UserData    *DeviceUserData `json:"user_data,omitempty"`
	Importance  *string         `json:"importance,omitempty"`
}

// Tag represents a custom tag
type Tag struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Colour string `json:"color"` // hex color code
}

// TagsResponse wraps the tags list from API
type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

// CreateTagRequest represents the request to create a new tag
type CreateTagRequest struct {
	Name   string `json:"name"`
	Colour string `json:"color"`
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	Name   *string `json:"name,omitempty"`
	Colour *string `json:"color,omitempty"`
}

// DeviceTagBinding represents the association between a device and a tag
type DeviceTagBinding struct {
	AgentID  int32 `json:"agent_id"`
	DeviceID int32 `json:"device_id"`
	TagID    int32 `json:"tag_id"`
}

// SNMPSensor represents an SNMP OID sensor
type SNMPSensor struct {
	ID        int32  `json:"id"`
	AgentID   int32  `json:"agent_id"`
	DeviceID  int32  `json:"device_id"`
	Name      string `json:"name"`
	OID       string `json:"oid"`
	Category  string `json:"category"`   // OTHER, etc.
	ValueType string `json:"value_type"` // STRING, NUMERIC, etc.
}

// CreateSNMPSensorRequest represents the request to create an SNMP sensor
type CreateSNMPSensorRequest struct {
	Name      string `json:"name"`
	OID       string `json:"oid"`
	Category  string `json:"category"`
	ValueType string `json:"value_type"`
}

// TCPSensor represents a TCP port sensor
type TCPSensor struct {
	ID       int32  `json:"id"`
	AgentID  int32  `json:"agent_id"`
	DeviceID int32  `json:"device_id"`
	Name     string `json:"name"`
	Port     int32  `json:"port"`
	Category string `json:"category"`
}

// CreateTCPSensorRequest represents the request to create a TCP sensor
type CreateTCPSensorRequest struct {
	Name     string `json:"name"`
	Port     int32  `json:"port"`
	Category string `json:"category"`
}

// Variable represents a device variable/metric
type Variable struct {
	ID            int32     `json:"id"`
	Label         string    `json:"label"`
	Path          string    `json:"path"`
	Value         string    `json:"value"`
	Unit          string    `json:"unit"`
	PreviousValue string    `json:"previous_value,omitempty"`
	Metric        string    `json:"metric"`
	UpdateTime    time.Time `json:"update_time"`
}

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	PageSize   int `json:"page_size,omitempty"`
	PageNumber int `json:"page_number,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
