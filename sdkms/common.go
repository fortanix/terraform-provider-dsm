// **********
// Terraform Provider - SDKMS: common functions
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1
//       - Date:      05/01/2021
// **********

package sdkms

import (
	"encoding/json"
	"fmt"
	//"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseHmg(d *schema.ResourceData, hmg map[string]interface{}) interface{} {
	var diags diag.Diagnostics

	rawHmg := make([]interface{}, len(hmg))
	counter := 0
	for rk, rv := range hmg {
		// FIXME: we know the structure within group api for hmg key
		if rawHmg[counter] == nil {
			rawHmg[counter] = make(map[string]interface{})
		}
		rawHmg[counter].(map[string]interface{})["hmg_id"] = rk
		for sk, sv := range rv.(map[string]interface{}) {
			if sk == "kind" || sk == "url" || sk == "access_key"{
				rawHmg[counter].(map[string]interface{})[sk] = sv
			} else if sk == "tls" {
				converted, err := json.Marshal(sv)
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Unable to call SDKMS provider API client",
						Detail:   fmt.Sprintf("[E]: COMMON: Parse HMG from Group: %s", err),
					})
					return diags
				}
				rawHmg[counter].(map[string]interface{})["tls"] = string(converted)
			}
		}
		counter++
    }

	return rawHmg
}