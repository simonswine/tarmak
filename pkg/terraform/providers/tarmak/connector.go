// Copyright Jetstack Ltd. See LICENSE for details.
package tarmak

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"

	tarmakRPC "github.com/jetstack/tarmak/pkg/terraform/providers/tarmak/rpc"
)

var rpcClient *rpc.Client
var rpcClientLock sync.Mutex

func RPCClient() (*rpc.Client, error) {
	connectorSocket := tarmakRPC.ConnectorSocket
	rpcClientLock.Lock()
	defer rpcClientLock.Unlock()
	if rpcClient == nil {
		tries := 20
		for {
			conn, err := rpc.Dial("unix", connectorSocket)
			rpcClient = conn
			if err == nil {
				return rpcClient, nil
			}
			if err != nil {
				fmt.Printf("unable to dial into unix socket '%s': %v\n", connectorSocket, err)
			}
			if tries == 0 {
				break
				return nil, fmt.Errorf("error connecting to connector socket '%s': %s", connectorSocket, err)
			}
			tries -= 1
			time.Sleep(time.Second)
		}
	}

	return rpcClient, nil
}
