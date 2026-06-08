package hue

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
)

const defaultDiscoveryEndpoint = "https://discovery.meethue.com/"

// Bridge describes a Hue bridge returned by the Philips Hue discovery service.
type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

// DiscoveryClient discovers Hue bridges through Philips Hue's N-UPnP service.
type DiscoveryClient struct {
	endpoint   string
	httpClient *http.Client
}

// NewDiscoveryClient creates a client for Philips Hue bridge discovery.
func NewDiscoveryClient() *DiscoveryClient {
	return &DiscoveryClient{
		endpoint: defaultDiscoveryEndpoint,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// DiscoverBridges returns Hue bridges announced through Philips Hue's discovery service.
func DiscoverBridges(ctx context.Context) ([]Bridge, error) {
	return NewDiscoveryClient().Discover(ctx)
}

// Discover returns Hue bridges announced through Philips Hue's discovery service.
func (client *DiscoveryClient) Discover(ctx context.Context) ([]Bridge, error) {
	endpoint := cmp.Or(client.endpoint, defaultDiscoveryEndpoint)
	httpClient := client.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create discovery request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get discovery service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery service returned %s", resp.Status)
	}

	var bridges []Bridge
	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, fmt.Errorf("decode discovery response: %w", err)
	}

	return normalizeDiscoveredBridges(bridges), nil
}

func normalizeDiscoveredBridges(bridges []Bridge) []Bridge {
	unique := make(map[string]Bridge, len(bridges))
	for _, bridge := range bridges {
		bridge.ID = strings.TrimSpace(bridge.ID)
		bridge.InternalIPAddress = strings.TrimSpace(bridge.InternalIPAddress)
		if bridge.InternalIPAddress == "" {
			continue
		}
		key := cmp.Or(bridge.ID, bridge.InternalIPAddress)
		unique[key] = bridge
	}

	normalized := make([]Bridge, 0, len(unique))
	for _, bridge := range unique {
		normalized = append(normalized, bridge)
	}

	slices.SortFunc(normalized, func(a, b Bridge) int {
		if idOrder := cmp.Compare(a.ID, b.ID); idOrder != 0 {
			return idOrder
		}
		return cmp.Compare(a.InternalIPAddress, b.InternalIPAddress)
	})

	return normalized
}
