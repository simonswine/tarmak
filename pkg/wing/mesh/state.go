package mesh

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/weaveworks/mesh"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
)

func NewWingStates(self mesh.PeerName) *WingStates {
	st := &WingStates{
		scheme: runtime.NewScheme(),
		instances: map[mesh.PeerName]tarmakv1alpha1.InstanceState{
			self: newInstanceState(self),
		},
	}
	st.codecs = serializer.NewCodecFactory(st.scheme)

	if err := tarmakv1alpha1.AddToScheme(st.scheme); err != nil {
		panic(err)
	}

	return st
}

func newInstanceState(self mesh.PeerName) tarmakv1alpha1.InstanceState {
	st := tarmakv1alpha1.InstanceState{
		Spec: &tarmakv1alpha1.InstanceStateSpec{},
		Status: &tarmakv1alpha1.InstanceStateStatus{
			Converge: &tarmakv1alpha1.InstanceStateStatusManifest{
				State: tarmakv1alpha1.InstanceStateConverging,
			},
		},
	}
	st.CreationTimestamp.Time = time.Now()
	st.Name = self.String()
	return st
}

type WingStates struct {
	lock   sync.Mutex
	self   mesh.PeerName
	scheme *runtime.Scheme
	codecs serializer.CodecFactory

	instances map[mesh.PeerName]tarmakv1alpha1.InstanceState
}

func (st *WingStates) Encode() [][]byte {
	var buf bytes.Buffer
	st.lock.Lock()
	defer st.lock.Unlock()
	return [][]byte{buf.Bytes()}
}

func (st *WingStates) encodeYAML(states *tarmakv1alpha1.InstanceStateList) ([]byte, error) {
	var encoder runtime.Encoder
	var buf bytes.Buffer

	mediaTypes := st.codecs.SupportedMediaTypes()
	for _, info := range mediaTypes {
		if info.MediaType == "application/yaml" {
			encoder = info.Serializer
			break
		}
	}
	if encoder == nil {
		return []byte{}, errors.New("unable to locate yaml encoder")
	}
	encoder = json.NewYAMLSerializer(json.DefaultMetaFactory, st.scheme, st.scheme)
	encoder = st.codecs.EncoderForVersion(encoder, tarmakv1alpha1.SchemeGroupVersion)

	if err := encoder.Encode(states, &buf); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil

}

// Merge merges the other GossipData into this one,
// and returns our resulting, complete state.
func (st *WingStates) Merge(other mesh.GossipData) (complete mesh.GossipData) {
	return st.mergeComplete(other.(*state).copy().set)
}

// Merge the set into our state, abiding increment-only semantics.
// Return a non-nil mesh.GossipData representation of the received set.
func (st *state) mergeReceived(set map[mesh.PeerName]int) (received mesh.GossipData) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	for peer, v := range set {
		if v <= st.set[peer] {
			delete(set, peer) // optimization: make the forwarded data smaller
			continue
		}
		st.set[peer] = v
	}

	return &state{
		set: set, // all remaining elements were novel to us
	}
}

// Merge the set into our state, abiding increment-only semantics.
// Return any key/values that have been mutated, or nil if nothing changed.
func (st *state) mergeDelta(set map[mesh.PeerName]int) (delta mesh.GossipData) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	for peer, v := range set {
		if v <= st.set[peer] {
			delete(set, peer) // requirement: it's not part of a delta
			continue
		}
		st.set[peer] = v
	}

	if len(set) <= 0 {
		return nil // per OnGossip requirements
	}
	return &state{
		set: set, // all remaining elements were novel to us
	}
}

// Merge the set into our state, abiding increment-only semantics.
// Return our resulting, complete state.
func (st *state) mergeComplete(set map[mesh.PeerName]int) (complete mesh.GossipData) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	for peer, v := range set {
		if v > st.set[peer] {
			st.set[peer] = v
		}
	}

	return &state{
		set: st.set, // n.b. can't .copy() due to lock contention
	}
}

	newInstances := other.(*WingStates).Instances
	for id, state := range newInstances {
		existing, ok := st.Instances[id]
		if !ok || existing.StatusChange.Before(state.StatusChange) {
			st.Instances[id] = state
		}
	}
	return st
}

// state implements GossipData.
var _ mesh.GossipData = &WingStates{}
