// Copyright Jetstack Ltd. See LICENSE for details.
package rpc

import (
	"io"
	"net/rpc"

	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
)

const (
	ConnectorSocket = "/tmp/tarmak-connector.sock"
)

type tarmakRPC struct {
	tarmak interfaces.Tarmak
}

// bind a new rpc server to socket
func Bind(tarmak interfaces.Tarmak, reader io.Reader, writer io.Writer, closer io.Closer) {

	tarmakRPC := tarmakRPC{tarmak: tarmak}

	s := rpc.NewServer()
	s.RegisterName("Tarmak", &tarmakRPC)

	tarmak.Log().Debugf("Connection made.")

	s.ServeConn(struct {
		io.Reader
		io.Writer
		io.Closer
	}{reader, writer, closer},
	)

}
