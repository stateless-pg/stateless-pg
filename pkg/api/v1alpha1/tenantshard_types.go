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
	TenantShardKind = "TenantShard"
	TenantShardKey  = "tenantshard"
	TenantShardName = "tenantshards"

	// ShardLayout versions for key->shard mapping algorithm
	ShardLayoutV1     int32 = 1
	ShardLayoutBroken int32 = 255

	// DefaultStripeSize is the default stripe size in pages (2048 pages = 16 MiB).
	// A lower stripe size distributes ingest load better across shards, but reduces IO amortization.
	// 16 MiB appears to be a reasonable balance: https://github.com/neondatabase/neon/pull/10510
	DefaultStripeSize int64 = 16 * 1024 / 8
)

// TenantShardId globally identifies a particular shard in a particular tenant.
//
// These are written as `<TenantId>-<ShardSlug>`, for example:
//   # The second shard in a two-shard tenant
//   072f1291a5310026820b2fe4b2968934-0102
//
// If the shard count is unsharded (1), the TenantShardId is written without
// a shard suffix and is equivalent to the encoding of a TenantId: this enables
// an unsharded TenantShardId to be used interchangably with a TenantId.
//
// The human-readable encoding of an unsharded TenantShardId, such as used in API URLs,
// is both forward and backward compatible with TenantId: a legacy TenantId can be
// decoded as a TenantShardId, and when re-encoded it will be parseable as a TenantId.
// +k8s:openapi-gen=true
type TenantShardId struct {
	// tenantId is the reference to the tenant this shard belongs to
	// +required
	TenantId v1.ObjectReference `json:"tenantId"`

	// shardNumber is the shard number within the shard count
	// Valid range: 0 to (shardCount - 1)
	// +required
	// +kubebuilder:validation:Minimum=0
	ShardNumber int32 `json:"shardNumber"`

	// shardCount is the total number of shards for the tenant
	// +required
	// +kubebuilder:validation:Minimum=1
	ShardCount int32 `json:"shardCount"`
}

// ShardIdentity uniquely identifies a shard and its configuration
// +k8s:openapi-gen=true
type ShardIdentity struct {
	// number defines the shard number within the shard count
	// Valid range: 0 to (count - 1)
	// +required
	// +kubebuilder:validation:Minimum=0
	Number int32 `json:"number"`

	// count defines the total number of shards for the tenant
	// This must match the shard count in all shards of the same tenant
	// +required
	// +kubebuilder:validation:Minimum=1
	Count int32 `json:"count"`

	// stripeSize defines the granularity in pages for distributing keys across shards
	// Default: 2048 pages (16 MiB)
	// A lower stripe size distributes ingest load better across shards, but reduces IO amortization.
	// +optional
	// +kubebuilder:default=2048
	// +kubebuilder:validation:Minimum=0
	StripeSize int64 `json:"stripeSize,omitempty"`

	// layout defines the layout version for key->shard mapping algorithm
	// This version number allows for future upgrades where the key->shard mapping may change
	// Valid values:
	// - 1: Current layout version (default)
	// - 255: Broken/unusable layout
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Enum=1;255
	Layout int32 `json:"layout,omitempty"`
}

// TenantShardSpec defines the desired state of TenantShard.
// +k8s:openapi-gen=true
type TenantShardSpec struct {
	// shardId globally identifies this shard
	// +required
	ShardId TenantShardId `json:"shardId"`

	// identity contains the core shard identification and configuration
	// +required
	Identity ShardIdentity `json:"identity"`

	// sequence is a runtime-only counter used to coordinate updates with background reconcilers
	// A reconciler runs to a particular sequence number to ensure consistency when multiple
	// reconcilers may be running concurrently
	// +optional
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	Sequence int64 `json:"sequence,omitempty"`
}

// TenantShardStatus defines the observed state of TenantShard.
// +k8s:openapi-gen=true
type TenantShardStatus struct {
	// conditions represent the current state of the TenantShard resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the shard is fully functional and serving requests
	// - "Progressing": the shard is being created, split, or updated
	// - "Degraded": the shard failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// phase represents the phase of the shard lifecycle
	// Valid values:
	// - "Creating": shard is being created
	// - "Attaching": shard is being attached to pageserver
	// - "Active": shard is ready to serve requests
	// - "Detaching": shard is being detached from pageserver
	// - "Deleted": shard has been deleted
	// +optional
	Phase string `json:"phase,omitempty"`

	// assignedPageserver is the name of the pageserver currently hosting this shard
	// +optional
	AssignedPageserver string `json:"assignedPageserver,omitempty"`

	// attachmentGeneration tracks the current attachment generation on the pageserver
	// Used to detect split-brain scenarios
	// +optional
	AttachmentGeneration *int64 `json:"attachmentGeneration,omitempty"`

	// lastUpdateTime is the last time this status was updated
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories="stateless-pg",shortName="ts"
// +kubebuilder:printcolumn:name="Tenant",type="string",JSONPath=".spec.tenantRef.name"
// +kubebuilder:printcolumn:name="Shard",type="string",JSONPath=".spec.shardIndex"
// +kubebuilder:printcolumn:name="Pageserver",type="string",JSONPath=".status.assignedPageserver"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.conditions[?(@.type == 'Available')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// TenantShard is the Schema for the tenantshards API
type TenantShard struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of TenantShard
	// +required
	Spec TenantShardSpec `json:"spec"`

	// status defines the observed state of TenantShard
	// +optional
	Status TenantShardStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// TenantShardList contains a list of TenantShard
// +k8s:openapi-gen=true
type TenantShardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []TenantShard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TenantShard{}, &TenantShardList{})
}
