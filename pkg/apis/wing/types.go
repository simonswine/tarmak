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

type InstanceSpec struct {
	ConvergeHash string
	DryRunPath   string
	DryRunHash   string
}

type InstanceStatus struct {
	Converge *InstanceStatusManifest
	DryRun   *InstanceStatusManifest
}

type InstanceStatusManifest struct {
	State string
	Hash  string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceList struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Items []Instance
}
