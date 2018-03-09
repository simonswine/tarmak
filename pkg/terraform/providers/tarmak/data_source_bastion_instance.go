// Copyright Jetstack Ltd. See LICENSE for details.
package tarmak

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/jetstack/tarmak/pkg/terraform/providers/tarmak/rpc"
)

func dataSourceBastionInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBastionInstanceRead,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBastionInstanceRead(d *schema.ResourceData, meta interface{}) error {

	client, err := RPCClient()
	if err != nil {
		return err
	}

	args := &rpc.BastionInstanceStatusArgs{
		Hostname: d.Get("hostname").(string),
		Username: d.Get("username").(string),
	}

	var reply rpc.BastionInstanceStatusReply
	err = client.Call("Tarmak.BastionInstanceStatus", args, &reply)
	if err != nil {
		return err
	}

	d.Set("status", reply.Status)

	return nil
}
