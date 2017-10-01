package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	InstanceID   string `json:"instanceID,omitempty"`
	InstancePool string `json:"instancePool,omitempty"`

	Spec   *InstanceSpec   `json:"spec,omitempty"`
	Status *InstanceStatus `json:"status,omitempty"`
}

type InstanceSpec struct {
	ConvergeHash string `json:"convergeHash,omitempty"`
	DryRunPath   string `json:"dryRunPath,omitempty"`
	DryRunHash   string `json:"dryRunHash,omitempty"`
}

type InstanceStatus struct {
	Converge *InstanceStatusManifest `json:"converge,omitempty"`
	DryRun   *InstanceStatusManifest `json:"dryRun,omitempty"`
}

type InstanceStatusManifest struct {
	State string `json:"state,omitempty"`
	Hash  string `json:"hash,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []Instance `json:"items"`
}
