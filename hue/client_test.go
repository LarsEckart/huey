package hue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api" {
			t.Errorf("expected /api, got %s", r.URL.Path)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["devicetype"] != "huey#test" {
			t.Errorf("expected devicetype huey#test, got %s", body["devicetype"])
		}

		w.Write([]byte(`[{"success":{"username":"abc123"}}]`))
	}))
	defer server.Close()

	// Extract host:port from server URL
	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "")

	username, err := client.Register("huey#test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if username != "abc123" {
		t.Errorf("expected username abc123, got %s", username)
	}
}

func TestRegister_LinkButtonNotPressed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"error":{"type":101,"address":"","description":"link button not pressed"}}]`))
	}))
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "")

	_, err := client.Register("huey#test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "link button not pressed") {
		t.Errorf("expected 'link button not pressed' error, got: %v", err)
	}
}

func TestGetLights(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/testuser/lights" {
			t.Errorf("expected /api/testuser/lights, got %s", r.URL.Path)
		}

		w.Write([]byte(`{
			"1": {
				"name": "Living Room",
				"type": "Extended color light",
				"state": {"on": true, "bri": 254, "hue": 10000, "sat": 200}
			},
			"2": {
				"name": "Bedroom",
				"type": "Dimmable light",
				"state": {"on": false, "bri": 100, "hue": 0, "sat": 0}
			}
		}`))
	}))
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "testuser")

	lights, err := client.GetLights()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lights) != 2 {
		t.Fatalf("expected 2 lights, got %d", len(lights))
	}

	// Find lights by ID (map iteration order is random)
	var light1, light2 *Light
	for i := range lights {
		if lights[i].ID == "1" {
			light1 = &lights[i]
		} else if lights[i].ID == "2" {
			light2 = &lights[i]
		}
	}

	if light1 == nil || light2 == nil {
		t.Fatal("missing expected lights")
	}

	if light1.Name != "Living Room" || !light1.On || light1.Brightness != 254 {
		t.Errorf("light1 data mismatch: %+v", light1)
	}
	if light2.Name != "Bedroom" || light2.On || light2.Brightness != 100 {
		t.Errorf("light2 data mismatch: %+v", light2)
	}
}

func TestGetLight(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/testuser/lights/1" {
			t.Errorf("expected /api/testuser/lights/1, got %s", r.URL.Path)
		}

		w.Write([]byte(`{
			"name": "Living Room",
			"type": "Extended color light",
			"state": {"on": true, "bri": 254, "hue": 10000, "sat": 200}
		}`))
	}))
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "testuser")

	light, err := client.GetLight("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if light.ID != "1" {
		t.Errorf("expected ID 1, got %s", light.ID)
	}
	if light.Name != "Living Room" {
		t.Errorf("expected name 'Living Room', got %s", light.Name)
	}
	if !light.On {
		t.Error("expected light to be on")
	}
}

func TestSetLightState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/testuser/lights/1/state" {
			t.Errorf("expected /api/testuser/lights/1/state, got %s", r.URL.Path)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if body["on"] != true {
			t.Errorf("expected on=true, got %v", body["on"])
		}

		w.Write([]byte(`[{"success":{"/lights/1/state/on":true}}]`))
	}))
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "testuser")

	on := true
	err := client.SetLightState("1", LightState{On: &on})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetLightState_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"error":{"type":3,"address":"/lights/999","description":"resource not available"}}]`))
	}))
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewClient(addr, "testuser")

	on := true
	err := client.SetLightState("999", LightState{On: &on})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "resource not available") {
		t.Errorf("expected 'resource not available' error, got: %v", err)
	}
}
