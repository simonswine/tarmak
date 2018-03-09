// Copyright Jetstack Ltd. See LICENSE for details.
package rpc

import (
	"fmt"
)

const (
	BastionInstanceStatusCall = "Tarmak.BastionInstanceStatus"
)

type BastionInstanceStatusArgs struct {
	Username string
	Hostname string
}

type BastionInstanceStatusReply struct {
	Status string
}

func (r *tarmakRPC) BastionInstanceStatus(args *BastionInstanceStatusArgs, result *BastionInstanceStatusReply) error {
	return fmt.Errorf("tarmak is not ready yet for your call")
}
