package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Client handles HTTP communication with a Hue bridge.
type Client struct {
	bridgeIP   string
	username   string
	httpClient *http.Client
}

// NewClient creates a Client for the given bridge IP and username.
// Username can be empty for registration calls.
func NewClient(bridgeIP, username string) *Client {
	return &Client{
		bridgeIP: bridgeIP,
		username: username,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// baseURL returns the API base URL.
func (c *Client) baseURL() string {
	return fmt.Sprintf("http://%s/api", c.bridgeIP)
}

// Register creates a new username on the bridge.
// Requires the bridge link button to be pressed first.
// deviceType format: "app_name#device_name" (e.g., "huey#macbook")
func (c *Client) Register(deviceType string) (string, error) {
	body := map[string]string{"devicetype": deviceType}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL(), "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	// Bridge returns an array of responses
	var results []map[string]any
	if err := json.Unmarshal(data, &results); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("empty response from bridge")
	}

	// Check for error
	if errData, ok := results[0]["error"].(map[string]any); ok {
		desc := errData["description"]
		return "", fmt.Errorf("bridge error: %v", desc)
	}

	// Extract username from success response
	if success, ok := results[0]["success"].(map[string]any); ok {
		if username, ok := success["username"].(string); ok {
			return username, nil
		}
	}

	return "", fmt.Errorf("unexpected response format: %s", string(data))
}

// Light represents a Hue light.
type Light struct {
	ID         string
	Name       string
	On         bool
	Brightness int // 0-254
	Hue        int // 0-65535
	Saturation int // 0-254
	Type       string
}

// lightResponse matches the JSON structure from the bridge for a single light.
type lightResponse struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	State struct {
		On         bool `json:"on"`
		Brightness int  `json:"bri"`
		Hue        int  `json:"hue"`
		Saturation int  `json:"sat"`
	} `json:"state"`
}

// GetLights returns all lights from the bridge.
func (c *Client) GetLights() ([]Light, error) {
	url := fmt.Sprintf("%s/%s/lights", c.baseURL(), c.username)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for error response (array format)
	if err := c.checkError(data); err != nil {
		return nil, err
	}

	// Bridge returns map of ID -> light object
	var lightsMap map[string]lightResponse
	if err := json.Unmarshal(data, &lightsMap); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	lights := make([]Light, 0, len(lightsMap))
	for id, lr := range lightsMap {
		lights = append(lights, Light{
			ID:         id,
			Name:       lr.Name,
			On:         lr.State.On,
			Brightness: lr.State.Brightness,
			Hue:        lr.State.Hue,
			Saturation: lr.State.Saturation,
			Type:       lr.Type,
		})
	}

	// Sort by ID numerically for natural order
	sort.Slice(lights, func(i, j int) bool {
		iID, _ := strconv.Atoi(lights[i].ID)
		jID, _ := strconv.Atoi(lights[j].ID)
		return iID < jID
	})

	return lights, nil
}

// GetLight returns a single light by ID.
func (c *Client) GetLight(id string) (*Light, error) {
	url := fmt.Sprintf("%s/%s/lights/%s", c.baseURL(), c.username, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if err := c.checkError(data); err != nil {
		return nil, err
	}

	var lr lightResponse
	if err := json.Unmarshal(data, &lr); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &Light{
		ID:         id,
		Name:       lr.Name,
		On:         lr.State.On,
		Brightness: lr.State.Brightness,
		Hue:        lr.State.Hue,
		Saturation: lr.State.Saturation,
		Type:       lr.Type,
	}, nil
}

// LightState represents the state to set on a light.
type LightState struct {
	On         *bool `json:"on,omitempty"`
	Brightness *int  `json:"bri,omitempty"`
	Hue        *int  `json:"hue,omitempty"`
	Saturation *int  `json:"sat,omitempty"`
}

// SetLightState changes the state of a light.
func (c *Client) SetLightState(id string, state LightState) error {
	url := fmt.Sprintf("%s/%s/lights/%s/state", c.baseURL(), c.username, id)

	jsonBody, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return c.checkError(data)
}

// RenameLight changes the name of a light.
func (c *Client) RenameLight(id string, name string) error {
	url := fmt.Sprintf("%s/%s/lights/%s", c.baseURL(), c.username, id)

	body := map[string]string{"name": name}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return c.checkError(data)
}

// checkError checks if the response contains an error.
// Bridge errors come as: [{"error":{"type":1,"address":"/...","description":"..."}}]
func (c *Client) checkError(data []byte) error {
	var results []map[string]any
	if err := json.Unmarshal(data, &results); err != nil {
		// Not an array response, so not an error format
		return nil
	}

	if len(results) == 0 {
		return nil
	}

	if errData, ok := results[0]["error"].(map[string]any); ok {
		desc := errData["description"]
		return fmt.Errorf("bridge error: %v", desc)
	}

	return nil
}

// Group represents a Hue group (room, zone, etc.).
type Group struct {
	ID      string
	Name    string
	Type    string   // "Room", "Zone", "LightGroup", etc.
	Lights  []string // Light IDs in this group
	AllOn   bool     // All lights in group are on
	AnyOn   bool     // At least one light is on
}

// groupResponse matches the JSON structure from the bridge for a single group.
type groupResponse struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Lights []string `json:"lights"`
	State  struct {
		AllOn bool `json:"all_on"`
		AnyOn bool `json:"any_on"`
	} `json:"state"`
}

// GetGroups returns all groups from the bridge.
func (c *Client) GetGroups() ([]Group, error) {
	url := fmt.Sprintf("%s/%s/groups", c.baseURL(), c.username)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if err := c.checkError(data); err != nil {
		return nil, err
	}

	// Bridge returns map of ID -> group object
	var groupsMap map[string]groupResponse
	if err := json.Unmarshal(data, &groupsMap); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	groups := make([]Group, 0, len(groupsMap))
	for id, gr := range groupsMap {
		groups = append(groups, Group{
			ID:     id,
			Name:   gr.Name,
			Type:   gr.Type,
			Lights: gr.Lights,
			AllOn:  gr.State.AllOn,
			AnyOn:  gr.State.AnyOn,
		})
	}

	// Sort by ID numerically
	sort.Slice(groups, func(i, j int) bool {
		iID, _ := strconv.Atoi(groups[i].ID)
		jID, _ := strconv.Atoi(groups[j].ID)
		return iID < jID
	})

	return groups, nil
}

// GroupAction represents the action to set on a group.
type GroupAction struct {
	On *bool `json:"on,omitempty"`
}

// SetGroupState changes the state of all lights in a group.
func (c *Client) SetGroupState(id string, action GroupAction) error {
	url := fmt.Sprintf("%s/%s/groups/%s/action", c.baseURL(), c.username, id)

	jsonBody, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return c.checkError(data)
}

// RenameGroup changes the name of a group.
func (c *Client) RenameGroup(id string, name string) error {
	url := fmt.Sprintf("%s/%s/groups/%s", c.baseURL(), c.username, id)

	body := map[string]string{"name": name}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return c.checkError(data)
}
