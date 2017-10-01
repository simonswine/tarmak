package mesh

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/weaveworks/mesh"
)

type Peer struct {
	states  *WingStates
	actions chan<- func()
	quit    chan struct{}
	send    mesh.Gossip
	log     *logrus.Entry
}

func (m *Mesh) newPeer() *Peer {
	p := &Peer{
		states:  NewWingStates(m.self),
		actions: make(chan func()),
		quit:    make(chan struct{}),
		log:     m.log,
	}
	return p
}

func (p *Peer) loop(actions <-chan func()) {
	for {
		select {
		case f := <-actions:
			f()
		case <-p.quit:
			return
		}
	}
}

// register the result of a mesh.Router.NewGossip.
func (p *Peer) register(send mesh.Gossip) {
	p.actions <- func() { p.send = send }
}

// Return a copy of our complete state.
func (p *Peer) Gossip() (complete mesh.GossipData) {
	return p.states
}

// Merge the gossiped data represented by buf into our state.
// Return the state information that was modified.
func (p *Peer) OnGossip(buf []byte) (delta mesh.GossipData, err error) {
	return nil, errors.New("unimplemented")
}
