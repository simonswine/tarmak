package connector

import (
	"io"
	"net"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Proxy struct {
	SocketPath string
	Done       chan struct{}
	log        *logrus.Entry
}

func NewProxy(socketPath string) *Proxy {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = os.Stderr

	return &Proxy{
		SocketPath: socketPath,
		Done:       make(chan struct{}),
		log: logger.WithFields(logrus.Fields{
			"socket": socketPath,
		}),
	}
}

func (p *Proxy) Start() error {
	p.log.Infoln("starting connector")
	listener, err := net.Listen("unix", p.SocketPath)
	if err != nil {
		return err
	}
	go p.run(listener)
	return nil
}

func (p *Proxy) Stop() {
	p.log.Infoln("stopping proxy")
	if p.Done == nil {
		return
	}
	close(p.Done)
	p.Done = nil
}

func (p *Proxy) run(listener net.Listener) {
	for {
		select {
		case <-p.Done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				p.handle(connection)
				p.Stop()
			} else {
				p.log.WithField("err", err).Errorln("Error accepting conn")
			}
		}
	}
}

func (p *Proxy) handle(connection net.Conn) {
	p.log.Debugln("Handling", connection)
	defer p.log.Debugln("Done handling", connection)
	defer connection.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(os.Stdout, connection, wg)
	go p.copy(connection, os.Stdin, wg)
	wg.Wait()
}

func (p *Proxy) copy(from io.Reader, to io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.Done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			p.log.WithField("err", err).Errorln("Error from copy")
			p.Stop()
			return
		}
	}
}
