package wing

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Instance struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	InstanceID   string
	InstancePool string

	Spec   *InstanceSpec
	Status *InstanceStatus
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	Converge *InstanceSpecManifest `json:"converge,omitempty"`
	DryRun   *InstanceSpecManifest `json:"dryRun,omitempty"`
}

//  InstaceSpecManifest defines location and hash for a specific manifest
type InstanceSpecManifest struct {
	Path             string      `json:"path,omitempty"`             // PATH to manifests (tar.gz)
	Hash             string      `json:"hash,omitempty"`             // md5 hash of manifests
	RequestTimestamp metav1.Time `json:"requestTimestamp,omitempty"` // timestamp when a converge was requested
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	Converge *InstanceStatusManifest `json:"converge,omitempty"`
	DryRun   *InstanceStatusManifest `json:"dryRun,omitempty"`
}

//  InstaceSpecManifest defines the state and hash of a run manifest
type InstanceStatusManifest struct {
	State               string      `json:"state,omitempty"`
	Hash                string      `json:"hash,omitempty"`                // md5 hash of manifests
	LastUpdateTimestamp metav1.Time `json:"lastUpdateTimestamp,omitempty"` // timestamp when a converge was requested
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceList struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Items []Instance
}
