package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_ListProviderRegions(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedCount  int
		expectedErr    bool
	}{
		{
			name: "successful response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/clusters/provider/regions"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiListProviderRegionsResp{
					Items: []*models.OpenapiListProviderRegionsItem{
						{
							CloudProvider: stringPtr("AWS"),
							Region:        stringPtr("us-west-2"),
							Available:     boolPtr(true),
						},
						{
							CloudProvider: stringPtr("AWS"),
							Region:        stringPtr("us-east-1"),
							Available:     boolPtr(true),
						},
						{
							CloudProvider: stringPtr("GCP"),
							Region:        stringPtr("us-central1"),
							Available:     boolPtr(true),
						},
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedCount: 3,
			expectedErr:   false,
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedCount: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test_public", "test_private")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			regions, err := client.ListProviderRegions()

			if tt.expectedErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(regions.Items) != tt.expectedCount {
				t.Errorf("Expected %d regions, got %d", tt.expectedCount, len(regions.Items))
			}

			// Verify first region if available
			if len(regions.Items) > 0 {
				firstRegion := regions.Items[0]
				if firstRegion.CloudProvider == nil || *firstRegion.CloudProvider == "" {
					t.Error("Expected cloud provider to be set")
				}
				if firstRegion.Region == nil || *firstRegion.Region == "" {
					t.Error("Expected region to be set")
				}
				if firstRegion.Available == nil {
					t.Error("Expected available status to be set")
				}
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}