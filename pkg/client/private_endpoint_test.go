package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

func TestClient_GetPrivateEndpointService(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResp   *models.OpenapiGetPrivateEndpointServiceResp
		expectedError  string
	}{
		{
			name:      "successful get private endpoint service",
			projectID: "test-project",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/clusters/test-cluster/private_endpoint_service" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/clusters/test-cluster/private_endpoint_service, got %s", r.URL.Path)
				}
				
				resp := &models.OpenapiGetPrivateEndpointServiceResp{
					CloudProvider: stringPtr("AWS"),
					Name:          stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
					Status:        stringPtr("ACTIVE"),
					DNSName:       stringPtr("vpce-svc-12345.us-east-1.vpce.amazonaws.com"),
					Port:          int64Ptr(4000),
					AzIDs:         []string{"use1-az1", "use1-az2"},
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectedResp: &models.OpenapiGetPrivateEndpointServiceResp{
				CloudProvider: stringPtr("AWS"),
				Name:          stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
				Status:        stringPtr("ACTIVE"),
				DNSName:       stringPtr("vpce-svc-12345.us-east-1.vpce.amazonaws.com"),
				Port:          int64Ptr(4000),
				AzIDs:         []string{"use1-az1", "use1-az2"},
			},
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
		{
			name:      "empty cluster ID",
			projectID: "test-project",
			clusterID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty cluster ID")
			},
			expectedError: "cluster ID is required",
		},
		{
			name:      "server error",
			projectID: "test-project",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    50000000,
					"message": "Internal server error",
				})
			},
			expectedError: "failed to execute request: TiDB Cloud API error (500): Internal server error (code: 50000000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.GetPrivateEndpointService(context.Background(), tt.projectID, tt.clusterID)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			if !privateEndpointServiceEqual(resp, tt.expectedResp) {
				t.Errorf("Expected response %+v, got %+v", tt.expectedResp, resp)
			}
		})
	}
}

func TestClient_CreatePrivateEndpointService(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResp   *models.OpenapiGetPrivateEndpointServiceResp
		expectedError  string
	}{
		{
			name:      "successful create private endpoint service",
			projectID: "test-project",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/clusters/test-cluster/private_endpoint_service" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/clusters/test-cluster/private_endpoint_service, got %s", r.URL.Path)
				}
				
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}
				
				resp := &models.OpenapiGetPrivateEndpointServiceResp{
					CloudProvider: stringPtr("AWS"),
					Name:          stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
					Status:        stringPtr("CREATING"),
					DNSName:       stringPtr("vpce-svc-12345.us-east-1.vpce.amazonaws.com"),
					Port:          int64Ptr(4000),
					AzIDs:         []string{"use1-az1", "use1-az2"},
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectedResp: &models.OpenapiGetPrivateEndpointServiceResp{
				CloudProvider: stringPtr("AWS"),
				Name:          stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
				Status:        stringPtr("CREATING"),
				DNSName:       stringPtr("vpce-svc-12345.us-east-1.vpce.amazonaws.com"),
				Port:          int64Ptr(4000),
				AzIDs:         []string{"use1-az1", "use1-az2"},
			},
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
		{
			name:      "empty cluster ID",
			projectID: "test-project",
			clusterID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty cluster ID")
			},
			expectedError: "cluster ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.CreatePrivateEndpointService(context.Background(), tt.projectID, tt.clusterID)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			if !privateEndpointServiceEqual(resp, tt.expectedResp) {
				t.Errorf("Expected response %+v, got %+v", tt.expectedResp, resp)
			}
		})
	}
}

func TestClient_ListPrivateEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResp   *models.OpenapiListPrivateEndpointsResp
		expectedError  string
	}{
		{
			name:      "successful list private endpoints",
			projectID: "test-project",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints, got %s", r.URL.Path)
				}
				
				resp := &models.OpenapiListPrivateEndpointsResp{
					Items: []*models.OpenapiPrivateEndpointItem{
						{
							ID:            stringPtr("pe-123"),
							CloudProvider: stringPtr("AWS"),
							ClusterID:     stringPtr("test-cluster"),
							Region:        stringPtr("us-east-1"),
							EndpointName:  stringPtr("vpce-12345"),
							Status:        stringPtr("ACTIVE"),
							Message:       stringPtr(""),
							ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
							ServiceStatus: stringPtr("ACTIVE"),
						},
					},
					Total: int64Ptr(1),
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectedResp: &models.OpenapiListPrivateEndpointsResp{
				Items: []*models.OpenapiPrivateEndpointItem{
					{
						ID:            stringPtr("pe-123"),
						CloudProvider: stringPtr("AWS"),
						ClusterID:     stringPtr("test-cluster"),
						Region:        stringPtr("us-east-1"),
						EndpointName:  stringPtr("vpce-12345"),
						Status:        stringPtr("ACTIVE"),
						Message:       stringPtr(""),
						ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
						ServiceStatus: stringPtr("ACTIVE"),
					},
				},
				Total: int64Ptr(1),
			},
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "test-cluster",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
		{
			name:      "empty cluster ID",
			projectID: "test-project",
			clusterID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty cluster ID")
			},
			expectedError: "cluster ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.ListPrivateEndpoints(context.Background(), tt.projectID, tt.clusterID)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			if !privateEndpointsEqual(resp, tt.expectedResp) {
				t.Errorf("Expected response %+v, got %+v", tt.expectedResp, resp)
			}
		})
	}
}

func TestClient_CreatePrivateEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		req            *models.OpenapiCreatePrivateEndpointReq
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResp   *models.OpenapiCreatePrivateEndpointResp
		expectedError  string
	}{
		{
			name:      "successful create private endpoint",
			projectID: "test-project",
			clusterID: "test-cluster",
			req: &models.OpenapiCreatePrivateEndpointReq{
				EndpointName: stringPtr("vpce-12345"),
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints, got %s", r.URL.Path)
				}
				
				var body models.OpenapiCreatePrivateEndpointReq
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}
				
				if body.EndpointName == nil || *body.EndpointName != "vpce-12345" {
					t.Errorf("Expected endpoint name 'vpce-12345', got %v", body.EndpointName)
				}
				
				resp := &models.OpenapiCreatePrivateEndpointResp{
					ID:            stringPtr("pe-123"),
					CloudProvider: stringPtr("AWS"),
					ClusterID:     stringPtr("test-cluster"),
					Region:        stringPtr("us-east-1"),
					EndpointName:  stringPtr("vpce-12345"),
					Status:        stringPtr("PENDING"),
					Message:       stringPtr("Creating private endpoint"),
					ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
					ServiceStatus: stringPtr("ACTIVE"),
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectedResp: &models.OpenapiCreatePrivateEndpointResp{
				ID:            stringPtr("pe-123"),
				CloudProvider: stringPtr("AWS"),
				ClusterID:     stringPtr("test-cluster"),
				Region:        stringPtr("us-east-1"),
				EndpointName:  stringPtr("vpce-12345"),
				Status:        stringPtr("PENDING"),
				Message:       stringPtr("Creating private endpoint"),
				ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
				ServiceStatus: stringPtr("ACTIVE"),
			},
		},
		{
			name:      "empty project ID",
			projectID: "",
			clusterID: "test-cluster",
			req: &models.OpenapiCreatePrivateEndpointReq{
				EndpointName: stringPtr("vpce-12345"),
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
		{
			name:      "empty cluster ID",
			projectID: "test-project",
			clusterID: "",
			req: &models.OpenapiCreatePrivateEndpointReq{
				EndpointName: stringPtr("vpce-12345"),
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty cluster ID")
			},
			expectedError: "cluster ID is required",
		},
		{
			name:      "nil request",
			projectID: "test-project",
			clusterID: "test-cluster",
			req:       nil,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with nil request")
			},
			expectedError: "request is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.CreatePrivateEndpoint(context.Background(), tt.projectID, tt.clusterID, tt.req)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			if !createPrivateEndpointEqual(resp, tt.expectedResp) {
				t.Errorf("Expected response %+v, got %+v", tt.expectedResp, resp)
			}
		})
	}
}

func TestClient_DeletePrivateEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		clusterID      string
		endpointID     string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  string
	}{
		{
			name:       "successful delete private endpoint",
			projectID:  "test-project",
			clusterID:  "test-cluster",
			endpointID: "pe-123",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Expected DELETE, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints/pe-123" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/clusters/test-cluster/private_endpoints/pe-123, got %s", r.URL.Path)
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{})
			},
		},
		{
			name:       "empty project ID",
			projectID:  "",
			clusterID:  "test-cluster",
			endpointID: "pe-123",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
		{
			name:       "empty cluster ID",
			projectID:  "test-project",
			clusterID:  "",
			endpointID: "pe-123",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty cluster ID")
			},
			expectedError: "cluster ID is required",
		},
		{
			name:       "empty endpoint ID",
			projectID:  "test-project",
			clusterID:  "test-cluster",
			endpointID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty endpoint ID")
			},
			expectedError: "endpoint ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			err = client.DeletePrivateEndpoint(context.Background(), tt.projectID, tt.clusterID, tt.endpointID)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestClient_ListPrivateEndpointsOfProject(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResp   *models.OpenapiListPrivateEndpointsResp
		expectedError  string
	}{
		{
			name:      "successful list private endpoints of project",
			projectID: "test-project",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1beta/projects/test-project/private_endpoints" {
					t.Errorf("Expected path /api/v1beta/projects/test-project/private_endpoints, got %s", r.URL.Path)
				}
				
				resp := &models.OpenapiListPrivateEndpointsResp{
					Items: []*models.OpenapiPrivateEndpointItem{
						{
							ID:            stringPtr("pe-123"),
							CloudProvider: stringPtr("AWS"),
							ClusterID:     stringPtr("cluster-1"),
							Region:        stringPtr("us-east-1"),
							EndpointName:  stringPtr("vpce-12345"),
							Status:        stringPtr("ACTIVE"),
							Message:       stringPtr(""),
							ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
							ServiceStatus: stringPtr("ACTIVE"),
						},
						{
							ID:            stringPtr("pe-456"),
							CloudProvider: stringPtr("AWS"),
							ClusterID:     stringPtr("cluster-2"),
							Region:        stringPtr("us-west-2"),
							EndpointName:  stringPtr("vpce-67890"),
							Status:        stringPtr("ACTIVE"),
							Message:       stringPtr(""),
							ServiceName:   stringPtr("com.amazonaws.vpce.us-west-2.vpce-svc-67890"),
							ServiceStatus: stringPtr("ACTIVE"),
						},
					},
					Total: int64Ptr(2),
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectedResp: &models.OpenapiListPrivateEndpointsResp{
				Items: []*models.OpenapiPrivateEndpointItem{
					{
						ID:            stringPtr("pe-123"),
						CloudProvider: stringPtr("AWS"),
						ClusterID:     stringPtr("cluster-1"),
						Region:        stringPtr("us-east-1"),
						EndpointName:  stringPtr("vpce-12345"),
						Status:        stringPtr("ACTIVE"),
						Message:       stringPtr(""),
						ServiceName:   stringPtr("com.amazonaws.vpce.us-east-1.vpce-svc-12345"),
						ServiceStatus: stringPtr("ACTIVE"),
					},
					{
						ID:            stringPtr("pe-456"),
						CloudProvider: stringPtr("AWS"),
						ClusterID:     stringPtr("cluster-2"),
						Region:        stringPtr("us-west-2"),
						EndpointName:  stringPtr("vpce-67890"),
						Status:        stringPtr("ACTIVE"),
						Message:       stringPtr(""),
						ServiceName:   stringPtr("com.amazonaws.vpce.us-west-2.vpce-svc-67890"),
						ServiceStatus: stringPtr("ACTIVE"),
					},
				},
				Total: int64Ptr(2),
			},
		},
		{
			name:      "empty project ID",
			projectID: "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				t.Error("Server should not be called with empty project ID")
			},
			expectedError: "project ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client, err := NewClient("test-key", "test-secret")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.ListPrivateEndpointsOfProject(context.Background(), tt.projectID)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			if !privateEndpointsEqual(resp, tt.expectedResp) {
				t.Errorf("Expected response %+v, got %+v", tt.expectedResp, resp)
			}
		})
	}
}

// Helper functions for comparison
func privateEndpointServiceEqual(a, b *models.OpenapiGetPrivateEndpointServiceResp) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	return stringPtrEqual(a.CloudProvider, b.CloudProvider) &&
		stringPtrEqual(a.Name, b.Name) &&
		stringPtrEqual(a.Status, b.Status) &&
		stringPtrEqual(a.DNSName, b.DNSName) &&
		int64PtrEqual(a.Port, b.Port) &&
		stringSliceEqual(a.AzIDs, b.AzIDs)
}

func privateEndpointsEqual(a, b *models.OpenapiListPrivateEndpointsResp) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	if !int64PtrEqual(a.Total, b.Total) {
		return false
	}
	
	if len(a.Items) != len(b.Items) {
		return false
	}
	
	for i, item := range a.Items {
		if !privateEndpointItemEqual(item, b.Items[i]) {
			return false
		}
	}
	
	return true
}

func privateEndpointItemEqual(a, b *models.OpenapiPrivateEndpointItem) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	return stringPtrEqual(a.ID, b.ID) &&
		stringPtrEqual(a.CloudProvider, b.CloudProvider) &&
		stringPtrEqual(a.ClusterID, b.ClusterID) &&
		stringPtrEqual(a.Region, b.Region) &&
		stringPtrEqual(a.EndpointName, b.EndpointName) &&
		stringPtrEqual(a.Status, b.Status) &&
		stringPtrEqual(a.Message, b.Message) &&
		stringPtrEqual(a.ServiceName, b.ServiceName) &&
		stringPtrEqual(a.ServiceStatus, b.ServiceStatus)
}

func createPrivateEndpointEqual(a, b *models.OpenapiCreatePrivateEndpointResp) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	return stringPtrEqual(a.ID, b.ID) &&
		stringPtrEqual(a.CloudProvider, b.CloudProvider) &&
		stringPtrEqual(a.ClusterID, b.ClusterID) &&
		stringPtrEqual(a.Region, b.Region) &&
		stringPtrEqual(a.EndpointName, b.EndpointName) &&
		stringPtrEqual(a.Status, b.Status) &&
		stringPtrEqual(a.Message, b.Message) &&
		stringPtrEqual(a.ServiceName, b.ServiceName) &&
		stringPtrEqual(a.ServiceStatus, b.ServiceStatus)
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Helper functions for pointer comparisons
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func int64PtrEqual(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}