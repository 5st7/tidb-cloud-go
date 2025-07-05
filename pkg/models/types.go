// Package models contains all request and response types for the TiDB Cloud API.
// These types are generated based on the OpenAPI specification and provide
// strong typing for all API operations including projects, clusters, backups,
// restores, and private endpoints.
package models

// Project API models
type OpenapiListProjectsResp struct {
	Items []*OpenapiListProjectItem `json:"items,omitempty"`
	Total *int64                    `json:"total,omitempty"`
}

type OpenapiListProjectItem struct {
	ID              *string `json:"id,omitempty"`
	OrgID           *string `json:"org_id,omitempty"`
	Name            *string `json:"name,omitempty"`
	ClusterCount    *int64  `json:"cluster_count,omitempty"`
	UserCount       *int64  `json:"user_count,omitempty"`
	CreateTimestamp *string `json:"create_timestamp,omitempty"`
}

type OpenapiCreateProjectReq struct {
	Name *string `json:"name,omitempty"`
}

type OpenapiCreateProjectResp struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// Provider Regions API models
type OpenapiListProviderRegionsResp struct {
	Items []*OpenapiListProviderRegionsItem `json:"items,omitempty"`
}

type OpenapiListProviderRegionsItem struct {
	CloudProvider *string `json:"cloud_provider,omitempty"`
	Region        *string `json:"region,omitempty"`
	Available     *bool   `json:"available,omitempty"`
}

// Cluster API models
type OpenapiListClustersOfProjectResp struct {
	Items []*OpenapiClusterItem `json:"items,omitempty"`
	Total *int64                `json:"total,omitempty"`
}

type OpenapiClusterItem struct {
	ID                *string                          `json:"id,omitempty"`
	Name              *string                          `json:"name,omitempty"`
	ClusterType       *string                          `json:"cluster_type,omitempty"`
	CloudProvider     *string                          `json:"cloud_provider,omitempty"`
	Region            *string                          `json:"region,omitempty"`
	Status            *OpenapiClusterItemStatus        `json:"status,omitempty"`
	Config            *OpenapiGetClusterConfig         `json:"config,omitempty"`
	ConnectionStrings *OpenapiClusterConnectionStrings `json:"connection_strings,omitempty"`
	CreateTimestamp   *string                          `json:"create_timestamp,omitempty"`
}

type OpenapiClusterItemStatus struct {
	ClusterStatus *string `json:"cluster_status,omitempty"`
}

type OpenapiGetClusterConfig struct {
	RootPassword *string                    `json:"root_password,omitempty"`
	Port         *int64                     `json:"port,omitempty"`
	Components   *OpenapiClusterComponents  `json:"components,omitempty"`
	IPAccessList []*OpenapiIpAccessListItem `json:"ip_access_list,omitempty"`
	Paused       *bool                      `json:"paused,omitempty"`
}

type OpenapiClusterComponents struct {
	TiDB    *OpenapiTiDBComponent    `json:"tidb,omitempty"`
	TiKV    *OpenapiTiKVComponent    `json:"tikv,omitempty"`
	TiFlash *OpenapiTiFlashComponent `json:"tiflash,omitempty"`
}

type OpenapiTiDBComponent struct {
	NodeSize     *string `json:"node_size,omitempty"`
	NodeQuantity *int64  `json:"node_quantity,omitempty"`
}

type OpenapiTiKVComponent struct {
	NodeSize       *string `json:"node_size,omitempty"`
	NodeQuantity   *int64  `json:"node_quantity,omitempty"`
	StorageSizeGib *int64  `json:"storage_size_gib,omitempty"`
}

type OpenapiTiFlashComponent struct {
	NodeSize       *string `json:"node_size,omitempty"`
	NodeQuantity   *int64  `json:"node_quantity,omitempty"`
	StorageSizeGib *int64  `json:"storage_size_gib,omitempty"`
}

type OpenapiClusterConnectionStrings struct {
	DefaultUser *string                    `json:"default_user,omitempty"`
	Standard    *OpenapiStandardConnection `json:"standard,omitempty"`
}

type OpenapiStandardConnection struct {
	Host *string `json:"host,omitempty"`
	Port *int64  `json:"port,omitempty"`
}

type OpenapiIpAccessListItem struct {
	CIDR        *string `json:"cidr,omitempty"`
	Description *string `json:"description,omitempty"`
}

type OpenapiCreateClusterReq struct {
	Name          *string               `json:"name,omitempty"`
	ClusterType   *string               `json:"cluster_type,omitempty"`
	CloudProvider *string               `json:"cloud_provider,omitempty"`
	Region        *string               `json:"region,omitempty"`
	Config        *OpenapiClusterConfig `json:"config,omitempty"`
}

type OpenapiClusterConfig struct {
	RootPassword *string                    `json:"root_password,omitempty"`
	Port         *int64                     `json:"port,omitempty"`
	Components   *OpenapiClusterComponents  `json:"components,omitempty"`
	IPAccessList []*OpenapiIpAccessListItem `json:"ip_access_list,omitempty"`
	Paused       *bool                      `json:"paused,omitempty"`
}

type OpenapiCreateClusterResp struct {
	ClusterID *string `json:"id,omitempty"`
}

type OpenapiUpdateClusterReq struct {
	Config *OpenapiUpdateClusterConfig `json:"config,omitempty"`
}

type OpenapiUpdateClusterConfig struct {
	Components *OpenapiUpdateClusterComponents `json:"components,omitempty"`
	Paused     *bool                           `json:"paused,omitempty"`
}

type OpenapiUpdateClusterComponents struct {
	TiDB    *OpenapiUpdateTiDBComponent    `json:"tidb,omitempty"`
	TiKV    *OpenapiUpdateTiKVComponent    `json:"tikv,omitempty"`
	TiFlash *OpenapiUpdateTiFlashComponent `json:"tiflash,omitempty"`
}

type OpenapiUpdateTiDBComponent struct {
	NodeSize     *string `json:"node_size,omitempty"`
	NodeQuantity *int64  `json:"node_quantity,omitempty"`
}

type OpenapiUpdateTiKVComponent struct {
	NodeSize       *string `json:"node_size,omitempty"`
	NodeQuantity   *int64  `json:"node_quantity,omitempty"`
	StorageSizeGib *int64  `json:"storage_size_gib,omitempty"`
}

type OpenapiUpdateTiFlashComponent struct {
	NodeSize       *string `json:"node_size,omitempty"`
	NodeQuantity   *int64  `json:"node_quantity,omitempty"`
	StorageSizeGib *int64  `json:"storage_size_gib,omitempty"`
}

// Backup API models
type OpenapiListBackupOfClusterResp struct {
	Items []*OpenapiListBackupItem `json:"items,omitempty"`
	Total *int64                   `json:"total,omitempty"`
}

type OpenapiListBackupItem struct {
	ID              *string                      `json:"id,omitempty"`
	Name            *string                      `json:"name,omitempty"`
	Description     *string                      `json:"description,omitempty"`
	ClusterID       *string                      `json:"cluster_id,omitempty"`
	Type            *string                      `json:"type,omitempty"`
	Status          *OpenapiListBackupItemStatus `json:"status,omitempty"`
	BackupTime      *string                      `json:"backup_time,omitempty"`
	ExpiryTime      *string                      `json:"expiry_time,omitempty"`
	BackupSizeBytes *int64                       `json:"backup_size_bytes,omitempty"`
	CreateTimestamp *string                      `json:"create_timestamp,omitempty"`
}

type OpenapiListBackupItemStatus struct {
	BackupStatus *string `json:"backup_status,omitempty"`
}

type OpenapiCreateBackupReq struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type OpenapiCreateBackupResp struct {
	BackupID *string `json:"backup_id,omitempty"`
}

type OpenapiGetBackupOfClusterResp struct {
	ID              *string                              `json:"id,omitempty"`
	Name            *string                              `json:"name,omitempty"`
	Description     *string                              `json:"description,omitempty"`
	ClusterID       *string                              `json:"cluster_id,omitempty"`
	Type            *string                              `json:"type,omitempty"`
	Status          *OpenapiGetBackupOfClusterRespStatus `json:"status,omitempty"`
	BackupTime      *string                              `json:"backup_time,omitempty"`
	ExpiryTime      *string                              `json:"expiry_time,omitempty"`
	BackupSizeBytes *int64                               `json:"backup_size_bytes,omitempty"`
	CreateTimestamp *string                              `json:"create_timestamp,omitempty"`
}

type OpenapiGetBackupOfClusterRespStatus struct {
	BackupStatus *string `json:"backup_status,omitempty"`
}

// Restore API models
type OpenapiListRestoreOfProjectResp struct {
	Items []*OpenapiListRestoreRespItem `json:"items,omitempty"`
	Total *int64                        `json:"total,omitempty"`
}

type OpenapiListRestoreRespItem struct {
	ID                *string                           `json:"id,omitempty"`
	Name              *string                           `json:"name,omitempty"`
	BackupID          *string                           `json:"backup_id,omitempty"`
	Status            *OpenapiListRestoreRespItemStatus `json:"status,omitempty"`
	ClusterInfo       *OpenapiClusterInfoOfRestore      `json:"cluster,omitempty"`
	CreateTimestamp   *string                           `json:"create_timestamp,omitempty"`
	FinishedTimestamp *string                           `json:"finished_timestamp,omitempty"`
}

type OpenapiListRestoreRespItemStatus struct {
	RestoreStatus *string `json:"restore_status,omitempty"`
}

type OpenapiClusterInfoOfRestore struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type OpenapiCreateRestoreReq struct {
	BackupID *string               `json:"backup_id,omitempty"`
	Name     *string               `json:"name,omitempty"`
	Config   *OpenapiClusterConfig `json:"config,omitempty"`
}

type OpenapiCreateRestoreResp struct {
	RestoreID *string `json:"restore_id,omitempty"`
}

type OpenapiGetRestoreResp struct {
	ID                *string                      `json:"id,omitempty"`
	Name              *string                      `json:"name,omitempty"`
	BackupID          *string                      `json:"backup_id,omitempty"`
	Status            *OpenapiGetRestoreRespStatus `json:"status,omitempty"`
	ClusterInfo       *OpenapiClusterInfoOfRestore `json:"cluster,omitempty"`
	CreateTimestamp   *string                      `json:"create_timestamp,omitempty"`
	FinishedTimestamp *string                      `json:"finished_timestamp,omitempty"`
}

type OpenapiGetRestoreRespStatus struct {
	RestoreStatus *string `json:"restore_status,omitempty"`
}

// Private Endpoint API models
type OpenapiGetPrivateEndpointServiceResp struct {
	CloudProvider *string  `json:"cloud_provider,omitempty"`
	Name          *string  `json:"name,omitempty"`
	Status        *string  `json:"status,omitempty"`
	DNSName       *string  `json:"dns_name,omitempty"`
	Port          *int64   `json:"port,omitempty"`
	AzIDs         []string `json:"az_ids,omitempty"`
}

type OpenapiListPrivateEndpointsResp struct {
	Items []*OpenapiPrivateEndpointItem `json:"items,omitempty"`
	Total *int64                        `json:"total,omitempty"`
}

type OpenapiPrivateEndpointItem struct {
	ID            *string `json:"id,omitempty"`
	CloudProvider *string `json:"cloud_provider,omitempty"`
	ClusterID     *string `json:"cluster_id,omitempty"`
	Region        *string `json:"region,omitempty"`
	EndpointName  *string `json:"endpoint_name,omitempty"`
	Status        *string `json:"status,omitempty"`
	Message       *string `json:"message,omitempty"`
	ServiceName   *string `json:"service_name,omitempty"`
	ServiceStatus *string `json:"service_status,omitempty"`
}

type OpenapiCreatePrivateEndpointReq struct {
	EndpointName *string `json:"endpoint_name,omitempty"`
}

type OpenapiCreatePrivateEndpointResp struct {
	ID            *string `json:"id,omitempty"`
	CloudProvider *string `json:"cloud_provider,omitempty"`
	ClusterID     *string `json:"cluster_id,omitempty"`
	Region        *string `json:"region,omitempty"`
	EndpointName  *string `json:"endpoint_name,omitempty"`
	Status        *string `json:"status,omitempty"`
	Message       *string `json:"message,omitempty"`
	ServiceName   *string `json:"service_name,omitempty"`
	ServiceStatus *string `json:"service_status,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Code    *int64        `json:"code,omitempty"`
	Message *string       `json:"message,omitempty"`
	Details []interface{} `json:"details,omitempty"`
}
