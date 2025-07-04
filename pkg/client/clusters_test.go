package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_ListClusters(t *testing.T) {
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
				expectedPath := "/api/v1beta/projects/project123/clusters"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiListClustersOfProjectResp{
					Items: []*models.OpenapiClusterItem{
						{
							ID:            stringPtr("cluster1"),
							Name:          stringPtr("Test Cluster"),
							ClusterType:   stringPtr("DEDICATED"),
							CloudProvider: stringPtr("AWS"),
							Region:        stringPtr("us-west-2"),
							Status: &models.OpenapiClusterItemStatus{
								ClusterStatus: stringPtr("AVAILABLE"),
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

			clusters, err := client.ListClusters(tt.projectID)

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

			if len(clusters.Items) != tt.expectedCount {
				t.Errorf("Expected %d clusters, got %d", tt.expectedCount, len(clusters.Items))
			}
		})
	}
}

func TestClient_GetCluster(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedName   string
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
				expectedPath := "/api/v1beta/projects/project123/clusters/cluster456"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				response := models.OpenapiClusterItem{
					ID:            stringPtr("cluster456"),
					Name:          stringPtr("My Test Cluster"),
					ClusterType:   stringPtr("DEDICATED"),
					CloudProvider: stringPtr("AWS"),
					Region:        stringPtr("us-west-2"),
					Status: &models.OpenapiClusterItemStatus{
						ClusterStatus: stringPtr("AVAILABLE"),
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedName: "My Test Cluster",
			expectedErr:  false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "cluster456",
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

			cluster, err := client.GetCluster(tt.projectID, tt.clusterID)

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

			if cluster.Name == nil || *cluster.Name != tt.expectedName {
				t.Errorf("Expected cluster name %s, got %v", tt.expectedName, cluster.Name)
			}
		})
	}
}

func TestClient_CreateCluster(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		request        *models.OpenapiCreateClusterReq
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedID     string
		expectedErr    bool
	}{
		{
			name:      "successful creation",
			projectID: "project123",
			request: &models.OpenapiCreateClusterReq{
				Name:          stringPtr("New Cluster"),
				ClusterType:   stringPtr("DEDICATED"),
				CloudProvider: stringPtr("AWS"),
				Region:        stringPtr("us-west-2"),
				Config: &models.OpenapiClusterConfig{
					RootPassword: stringPtr("testpassword"),
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
				expectedPath := "/api/v1beta/projects/project123/clusters"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s, got %s", expectedPath, r.URL.Path)
				}

				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="tidbcloud", nonce="test123", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Verify request body
				var reqBody models.OpenapiCreateClusterReq
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}

				response := models.OpenapiCreateClusterResp{
					ClusterID: stringPtr("new-cluster-123"),
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedID:  "new-cluster-123",
			expectedErr: false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			request:   &models.OpenapiCreateClusterReq{},
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

			response, err := client.CreateCluster(tt.projectID, tt.request)

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

			if response.ClusterID == nil || *response.ClusterID != tt.expectedID {
				t.Errorf("Expected cluster ID %s, got %v", tt.expectedID, response.ClusterID)
			}
		})
	}
}