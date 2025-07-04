package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_ListBackups(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedCount  int
		expectedErr    bool
	}{
		{
			name:      "successful response",
			projectID: "project123",
			clusterID: "cluster456",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/projects/project123/clusters/cluster456/backups"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiListBackupOfClusterResp{
					Items: []*models.OpenapiListBackupItem{
						{
							ID:          stringPtr("backup1"),
							Name:        stringPtr("Daily Backup"),
							Description: stringPtr("Automated daily backup"),
							ClusterID:   stringPtr("cluster456"),
							Type:        stringPtr("MANUAL"),
							Status: &models.OpenapiListBackupItemStatus{
								BackupStatus: stringPtr("SUCCESS"),
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
			clusterID: "cluster456",
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

			backups, err := client.ListBackups(tt.projectID, tt.clusterID)

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

			if len(backups.Items) != tt.expectedCount {
				t.Errorf("Expected %d backups, got %d", tt.expectedCount, len(backups.Items))
			}
		})
	}
}

func TestClient_CreateBackup(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		request        *models.OpenapiCreateBackupReq
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedID     string
		expectedErr    bool
	}{
		{
			name:      "successful creation",
			projectID: "project123",
			clusterID: "cluster456",
			request: &models.OpenapiCreateBackupReq{
				Name:        stringPtr("Test Backup"),
				Description: stringPtr("Test backup description"),
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1beta/projects/project123/clusters/cluster456/backups"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiCreateBackupResp{
					BackupID: stringPtr("new-backup-123"),
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedID:  "new-backup-123",
			expectedErr: false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "cluster456",
			request:   &models.OpenapiCreateBackupReq{},
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

			response, err := client.CreateBackup(tt.projectID, tt.clusterID, tt.request)

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

			if response.BackupID == nil || *response.BackupID != tt.expectedID {
				t.Errorf("Expected backup ID %s, got %v", tt.expectedID, response.BackupID)
			}
		})
	}
}