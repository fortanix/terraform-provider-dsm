// **********
// Terraform Provider - DSM: main provider program
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.2.4
//       - Date:      27/07/2021
// **********

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-dsm/dsm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return dsm.Provider()
		},
	})
}
