/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TenantKind = "Tenant"
	TenantKey  = "tenant"
	TenantName = "tenants"
)

// TenantSpec defines the desired state of Tenant.
// +k8s:openapi-gen=true
type TenantSpec struct {
	
	// generation defines version number for split-brain safety
	// +optional
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	Generation int64 `json:"generation,omitempty"`

	// placementPolicy defines the placement strategy for the tenant
	// Valid values:
	// - Attached(0): Single node, development (default)
	// - Attached(1): HA with failover
	// - Attached(2): HA with multiple replicas
	// - Secondary: Onboarding/migration standby
	// - Detached: Archived/idle data
	// +optional
	// +kubebuilder:default="Attached(0)"
	// +kubebuilder:validation:Enum="Attached(0)";"Attached(1)";"Attached(2)";Secondary;Detached
	PlacementPolicy string `json:"placementPolicy,omitempty"`

	// config defines the configuration parameters for the tenant
	// +optional
	Config TenantConfig `json:"config,omitempty"`
}

// TenantConfig defines the desired state of Tenant.
// +k8s:openapi-gen=true
type TenantConfig struct {
	// neonClusterRef is a reference to the NeonCluster resource this tenant belongs to
	// +required
	NeonClusterRef *v1.ObjectReference `json:"neonClusterRef"`

	// checkpointDistance defines the size threshold between checkpoints (L0 layer file size)
	// Default: 268435456 (256 MB)
	// +optional
	// +kubebuilder:default=268435456
	// +kubebuilder:validation:Minimum=0
	CheckpointDistance int64 `json:"checkpointDistance,omitempty"`

	// checkpointTimeout defines the maximum time between checkpoints
	// Default: 10 minutes
	// +optional
	// +kubebuilder:default="10m"
	CheckpointTimeout string `json:"checkpointTimeout,omitempty"`

	// compactionTargetSize defines the target layer size for image/delta layers (L1 layer file size)
	// Default: 134217728 (128 MB)
	// +optional
	// +kubebuilder:default=134217728
	// +kubebuilder:validation:Minimum=0
	CompactionTargetSize int64 `json:"compactionTargetSize,omitempty"`

	// compactionPeriod defines how often to run compaction
	// Default: 20 seconds
	// +optional
	// +kubebuilder:default="20s"
	CompactionPeriod string `json:"compactionPeriod,omitempty"`

	// compactionThreshold defines the threshold for triggering compaction
	// Default: 10 layers
	// +optional
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=0
	CompactionThreshold int32 `json:"compactionThreshold,omitempty"`

	// compactionUpperLimit defines the maximum layers to compact before memory limits
	// Default: 10 layers
	// +optional
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=0
	CompactionUpperLimit int32 `json:"compactionUpperLimit,omitempty"`

	// compactionAlgorithm defines the algorithm to use for layer compaction
	// +optional
	// +kubebuilder:default="Legacy"
	CompactionAlgorithm string `json:"compactionAlgorithm,omitempty"`

	// compactionShardAncestor enables shard-aware compaction
	// +optional
	// +kubebuilder:default=true
	CompactionShardAncestor bool `json:"compactionShardAncestor,omitempty"`

	// compactionL0First enables L0 compaction pass for responsiveness
	// +optional
	// +kubebuilder:default=true
	CompactionL0First bool `json:"compactionL0First,omitempty"`

	// compactionL0Semaphore enables L0 compaction semaphore for concurrency control
	// +optional
	// +kubebuilder:default=true
	CompactionL0Semaphore bool `json:"compactionL0Semaphore,omitempty"`

	// l0FlushDelayThreshold defines the L0 flush delay threshold
	// +optional
	// +kubebuilder:validation:Minimum=0
	L0FlushDelayThreshold *int32 `json:"l0FlushDelayThreshold,omitempty"`

	// l0FlushStallThreshold defines the L0 stall threshold
	// +optional
	// +kubebuilder:validation:Minimum=0
	L0FlushStallThreshold *int32 `json:"l0FlushStallThreshold,omitempty"`

	// gcHorizon defines how far back to retain data
	// Default: 67108864 (64 MB)
	// +optional
	// +kubebuilder:default=67108864
	// +kubebuilder:validation:Minimum=0
	GcHorizon int64 `json:"gcHorizon,omitempty"`

	// gcPeriod defines how often to run garbage collection
	// Default: 1 hour
	// +optional
	// +kubebuilder:default="1h"
	GcPeriod string `json:"gcPeriod,omitempty"`

	// imageCreationThreshold defines after how many L0 layers to create an image layer
	// Default: 3 layers
	// +optional
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=0
	ImageCreationThreshold int32 `json:"imageCreationThreshold,omitempty"`

	// imageLayerForceCreationPeriod forces image layer creation at interval
	// +optional
	ImageLayerForceCreationPeriod *string `json:"imageLayerForceCreationPeriod,omitempty"`

	// pitrInterval defines the point-in-time recovery window
	// Default: 7 days retention
	// +optional
	// +kubebuilder:default="168h"
	PitrInterval string `json:"pitrInterval,omitempty"`

	// walreceiverConnectTimeout defines the timeout for WAL receiver connection
	// Default: 10 seconds
	// +optional
	// +kubebuilder:default="10s"
	WalreceiverConnectTimeout string `json:"walreceiverConnectTimeout,omitempty"`

	// laggingWalTimeout defines the timeout before disconnecting slow WAL receiver
	// Default: 10 seconds
	// +optional
	// +kubebuilder:default="10s"
	LaggingWalTimeout string `json:"laggingWalTimeout,omitempty"`

	// maxLsnWalLag defines the maximum WAL lag allowed (supports 1GB/s throughput)
	// Default: 1073741824 (1 GB)
	// +optional
	// +kubebuilder:default=1073741824
	// +kubebuilder:validation:Minimum=0
	MaxLsnWalLag *int64 `json:"maxLsnWalLag,omitempty"`

	// evictionPolicy defines the page eviction strategy
	// Default: NoEviction (don't evict pages)
	// Valid values: NoEviction, AlwaysEvict
	// +optional
	// +kubebuilder:default="NoEviction"
	// +kubebuilder:validation:Enum=NoEviction;AlwaysEvict
	EvictionPolicy string `json:"evictionPolicy,omitempty"`

	// minResidentSizeOverride overrides the minimum resident size threshold
	// +optional
	// +kubebuilder:validation:Minimum=0
	MinResidentSizeOverride *int64 `json:"minResidentSizeOverride,omitempty"`

	// evictionsLowResidenceDurationMetricThreshold defines the metric threshold for low residence time
	// +optional
	// +kubebuilder:default="24h"
	EvictionsLowResidenceDurationMetricThreshold string `json:"evictionsLowResidenceDurationMetricThreshold,omitempty"`

	// heatmapPeriod defines the heatmap refresh period (0 = disabled)
	// +optional
	// +kubebuilder:default="0s"
	HeatmapPeriod string `json:"heatmapPeriod,omitempty"`

	// lazySlruDownload enables lazy download strategy
	// +optional
	// +kubebuilder:default=false
	LazySlruDownload bool `json:"lazySlruDownload,omitempty"`

	// timelineGetThrottle defines the rate limiting configuration for page reads
	// +optional
	TimelineGetThrottle *ThrottleConfig `json:"timelineGetThrottle,omitempty"`

	// imageLayerCreationCheckThreshold defines how many L0 layers to ingest WAL for before checking image creation
	// +optional
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=0
	ImageLayerCreationCheckThreshold int32 `json:"imageLayerCreationCheckThreshold,omitempty"`

	// imageCreationPreemptThreshold preempts image creation if L0 backpressure
	// +optional
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=0
	ImageCreationPreemptThreshold int32 `json:"imageCreationPreemptThreshold,omitempty"`

	// lsnLeaseLength defines the LSN lease duration for read consistency
	// Default: 10 minutes
	// +optional
	// +kubebuilder:default="10m"
	LsnLeaseLength string `json:"lsnLeaseLength,omitempty"`

	// lsnLeaseLengthForTs defines the LSN lease duration for timestamp queries
	// Default: 60 seconds
	// +optional
	// +kubebuilder:default="60s"
	LsnLeaseLengthForTs string `json:"lsnLeaseLengthForTs,omitempty"`

	// timelineOffloading enables timeline offloading to secondary nodes
	// Default: true (enabled)
	// +optional
	// +kubebuilder:default=true
	TimelineOffloading bool `json:"timelineOffloading,omitempty"`

	// relSizeV2Enabled enables relation size v2 calculation
	// +optional
	// +kubebuilder:default=false
	RelSizeV2Enabled bool `json:"relSizeV2Enabled,omitempty"`

	// gcCompactionEnabled enables GC compaction pass
	// Default: true (enabled)
	// +optional
	// +kubebuilder:default=true
	GcCompactionEnabled bool `json:"gcCompactionEnabled,omitempty"`

	// gcCompactionVerification verifies GC compaction results
	// Default: true (enabled)
	// +optional
	// +kubebuilder:default=true
	GcCompactionVerification bool `json:"gcCompactionVerification,omitempty"`

	// gcCompactionInitialThresholdKb defines the initial threshold for GC compaction in KB
	// Default: 5368709 (5 GB)
	// +optional
	// +kubebuilder:default=5368709
	// +kubebuilder:validation:Minimum=0
	GcCompactionInitialThresholdKb int64 `json:"gcCompactionInitialThresholdKb,omitempty"`

	// gcCompactionRatioPercent defines the compaction ratio percentage for GC
	// Default: 100%
	// +optional
	// +kubebuilder:default=100
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	GcCompactionRatioPercent int64 `json:"gcCompactionRatioPercent,omitempty"`

	// samplingRatio defines the sampling configuration for metrics
	// +optional
	SamplingRatio *SamplingRatio `json:"samplingRatio,omitempty"`

	// relsizeSnapshotCacheCapacity defines the snapshot cache size for relation size
	// Default: 1000 entries
	// +optional
	// +kubebuilder:default=1000
	// +kubebuilder:validation:Minimum=0
	RelsizeSnapshotCacheCapacity int32 `json:"relsizeSnapshotCacheCapacity,omitempty"`

	// basebackupCacheEnabled enables basebackup caching
	// +optional
	// +kubebuilder:default=false
	BasebackupCacheEnabled bool `json:"basebackupCacheEnabled,omitempty"`
}

// ShardParameters defines sharding configuration for the tenant
// +k8s:openapi-gen=true
type ShardParameters struct {
	// count defines the number of shards (1 = unsharded)
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Count int32 `json:"count,omitempty"`

	// stripeSize defines the granularity (in bytes) for distributing keys across shards
	// Common values: 8192 (8 KB - Default), 16384 (16 KB), 32768 (32 KB)
	// +optional
	// +kubebuilder:default=8192
	// +kubebuilder:validation:Minimum=0
	StripeSize int64 `json:"stripe_size,omitempty"`
}

// ThrottleConfig defines rate limiting configuration for page reads
// +k8s:openapi-gen=true
type ThrottleConfig struct {
	// enabled defines if throttling is enabled
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// rps defines requests per second limit
	// +optional
	// +kubebuilder:validation:Minimum=0
	Rps *int32 `json:"rps,omitempty"`
}

// SamplingRatio defines sampling configuration for metrics
// +k8s:openapi-gen=true
type SamplingRatio struct {
	// ratio defines the sampling ratio value
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	Ratio *string `json:"ratio,omitempty"`
}

// TenantStatus defines the observed state of Tenant.
// +k8s:openapi-gen=true
type TenantStatus struct {
	// conditions represent the current state of the Tenant resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// phase represents the phase of the tenant creation/deletion
	// +optional
	Phase string `json:"phase,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories="stateless-pg",shortName="tn"
// +kubebuilder:printcolumn:name="Database",type="string",JSONPath=".spec.databaseName"
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.conditions[?(@.type == 'Available')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// Tenant is the Schema for the tenants API
type Tenant struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Tenant
	// +required
	Spec TenantSpec `json:"spec"`

	// status defines the observed state of Tenant
	// +optional
	Status TenantStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant
// +k8s:openapi-gen=true
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
