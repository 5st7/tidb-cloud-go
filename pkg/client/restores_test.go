package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_ListRestores(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedCount  int
		expectedErr    bool
	}{
		{
			name:      "successful response",
			projectID: "project123",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/projects/project123/restores"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiListRestoreOfProjectResp{
					Items: []*models.OpenapiListRestoreRespItem{
						{
							ID:       stringPtr("restore1"),
							Name:     stringPtr("Test Restore"),
							BackupID: stringPtr("backup123"),
							Status: &models.OpenapiListRestoreRespItemStatus{
								RestoreStatus: stringPtr("SUCCESS"),
							},
							ClusterInfo: &models.OpenapiClusterInfoOfRestore{
								ID:   stringPtr("new-cluster-789"),
								Name: stringPtr("Restored Cluster"),
							},
						},
					},
					Total: int64Ptr(1),
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Should not be called
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

			restores, err := client.ListRestores(tt.projectID)

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

			if len(restores.Items) != tt.expectedCount {
				t.Errorf("Expected %d restores, got %d", tt.expectedCount, len(restores.Items))
			}
		})
	}
}

func TestClient_CreateRestore(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		request        *models.OpenapiCreateRestoreReq
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedID     string
		expectedErr    bool
	}{
		{
			name:      "successful creation",
			projectID: "project123",
			request: &models.OpenapiCreateRestoreReq{
				BackupID: stringPtr("backup123"),
				Name:     stringPtr("Test Restore"),
				Config: &models.OpenapiClusterConfig{
					RootPassword: stringPtr("newpassword"),
					Port:         int64Ptr(4000),
					Components: &models.OpenapiClusterComponents{
						TiDB: &models.OpenapiTiDBComponent{
							NodeSize:     stringPtr("8C16G"),
							NodeQuantity: int64Ptr(1),
						},
						TiKV: &models.OpenapiTiKVComponent{
							NodeSize:       stringPtr("8C32G"),
							NodeQuantity:   int64Ptr(3),
							StorageSizeGib: int64Ptr(500),
						},
					},
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/projects/project123/restores"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiCreateRestoreResp{
					RestoreID: stringPtr("new-restore-123"),
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedID:  "new-restore-123",
			expectedErr: false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			request:   &models.OpenapiCreateRestoreReq{},
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

			response, err := client.CreateRestore(tt.projectID, tt.request)

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

			if response.RestoreID == nil || *response.RestoreID != tt.expectedID {
				t.Errorf("Expected restore ID %s, got %v", tt.expectedID, response.RestoreID)
			}
		})
	}
}
