package hue

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDiscoveryClientDiscover(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		_, _ = w.Write([]byte(`[
			{"id":"bridge-b","internalipaddress":"192.168.1.20"},
			{"id":"bridge-a","internalipaddress":"192.168.1.10"},
			{"id":"empty-ip","internalipaddress":""},
			{"id":" bridge-a ","internalipaddress":" 192.168.1.10 "}
		]`))
	}))
	defer server.Close()

	client := &DiscoveryClient{
		endpoint:   server.URL,
		httpClient: server.Client(),
	}

	bridges, err := client.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(bridges) != 2 {
		t.Fatalf("expected 2 bridges, got %d: %+v", len(bridges), bridges)
	}

	if bridges[0].ID != "bridge-a" || bridges[0].InternalIPAddress != "192.168.1.10" {
		t.Errorf("first bridge mismatch: %+v", bridges[0])
	}
	if bridges[1].ID != "bridge-b" || bridges[1].InternalIPAddress != "192.168.1.20" {
		t.Errorf("second bridge mismatch: %+v", bridges[1])
	}
}

func TestDiscoveryClientDiscover_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := &DiscoveryClient{
		endpoint:   server.URL,
		httpClient: server.Client(),
	}

	_, err := client.Discover(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDiscoveryClientDiscover_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	client := &DiscoveryClient{
		endpoint:   server.URL,
		httpClient: server.Client(),
	}

	_, err := client.Discover(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
