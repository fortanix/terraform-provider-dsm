// **********
// Terraform Provider - DSM: common functions
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.2
//       - Date:      05/01/2021
// **********

package dsm

import (
	//"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	//"encoding/pem"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/crypto/ssh"
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
			if sk == "kind" || sk == "url" || sk == "access_key" {
				rawHmg[counter].(map[string]interface{})[sk] = sv
			} else if sk == "tls" {
				converted, err := json.Marshal(sv)
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[DSM SDK] Unable to call DSM provider API client",
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

func PublicPEMtoOpenSSH(pemBytes []byte) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	derBlock, err := base64.StdEncoding.DecodeString(string(pemBytes))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[GENERAL SDK] DER Conversion Error",
			Detail:   fmt.Sprintf("[E]: COMMON: Unable convert PEM public key to DER: %s", err),
		})
		return nil, diags
	}

	pubKey, err := x509.ParsePKIXPublicKey(derBlock)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[GENERAL SDK] PEM Conversion Error",
			Detail:   fmt.Sprintf("[E]: COMMON: Unable generate X.509 parsed public key: %s", err),
		})
		return nil, diags
	}

	pub, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[GENERAL SDK] PEM Conversion Error",
			Detail:   fmt.Sprintf("[E]: COMMON: Unable to create OpenSSH format public key: %s", err),
		})
		return nil, diags
	}

	sshPubKey := base64.StdEncoding.EncodeToString(pub.Marshal())

	return &sshPubKey, diags
}
