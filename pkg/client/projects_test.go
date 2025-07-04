package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_ListProjects(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedProjects int
		expectedErr    bool
	}{
		{
			name: "successful response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects" {
					t.Errorf("Expected /api/v1beta/projects, got %s", r.URL.Path)
				}

				// First request without auth header should return 401
				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", qop="auth", opaque="5ccc069c403ebaf9f0171e9517f40e41", algorithm="MD5"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Return successful response
				response := models.OpenapiListProjectsResp{
					Items: []*models.OpenapiListProjectItem{
						{
							ID:           stringPtr("project1"),
							Name:         stringPtr("Test Project"),
							ClusterCount: int64Ptr(1),
							UserCount:    int64Ptr(1),
						},
					},
					Total: int64Ptr(1),
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedProjects: 1,
			expectedErr:      false,
		},
		{
			name: "unauthorized without retry",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    49900001,
					"message": "public_key not found",
				})
			},
			expectedProjects: 0,
			expectedErr:      true,
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

			projects, err := client.ListProjects()

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

			if len(projects.Items) != tt.expectedProjects {
				t.Errorf("Expected %d projects, got %d", tt.expectedProjects, len(projects.Items))
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestClient_CreateProject(t *testing.T) {
	tests := []struct {
		name           string
		request        *models.OpenapiCreateProjectReq
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedID     string
		expectedErr    bool
	}{
		{
			name: "successful creation",
			request: &models.OpenapiCreateProjectReq{
				Name: stringPtr("New Project"),
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/projects"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Verify request body
				var reqBody models.OpenapiCreateProjectReq
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}

				response := models.OpenapiCreateProjectResp{
					ID:   stringPtr("new-project-123"),
					Name: stringPtr("New Project"),
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedID:  "new-project-123",
			expectedErr: false,
		},
		{
			name:    "nil request",
			request: nil,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Should not be called
			},
			expectedErr: true,
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

			response, err := client.CreateProject(tt.request)

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

			if response.ID == nil || *response.ID != tt.expectedID {
				t.Errorf("Expected project ID %s, got %v", tt.expectedID, response.ID)
			}
		})
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}