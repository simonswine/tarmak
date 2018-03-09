// Copyright Jetstack Ltd. See LICENSE for details.
package tarmak

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// TODO: Move the validation to this, requires conditional schemas
	// TODO: Move the configuration to this, requires validation

	// The actual provider
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"tarmak_bastion_instance": dataSourceBastionInstance(),
		},
	}

}
