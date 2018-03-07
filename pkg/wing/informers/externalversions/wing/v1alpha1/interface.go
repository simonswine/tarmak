// Copyright Jetstack Ltd. See LICENSE for details.

// This file was automatically generated by informer-gen

package v1alpha1

import (
	internalinterfaces "github.com/jetstack/tarmak/pkg/wing/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Instances returns a InstanceInformer.
	Instances() InstanceInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Instances returns a InstanceInformer.
func (v *version) Instances() InstanceInformer {
	return &instanceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
