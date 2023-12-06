// **********
// Terraform Provider - DSM: common functions
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
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
	"strconv"

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

// computes the add and delete arrays - Ravi Gopal
func compute_add_and_del_arrays(old_array interface{}, new_array interface{}) ([]string, []string) {

	/*
	* Compares old state and new state
	* segregates the arrays to be added and arrays to be deleted.
	*/
	old_array_set := old_array.([]interface{})
	new_array_set := new_array.([]interface{})

	old_array_ids := make([]string, len(old_array_set))
	for i, v := range old_array_set {
		old_array_ids[i] = v.(string)
	}
	new_array_ids := make([]string, len(new_array_set))
	for i, v := range new_array_set {
		new_array_ids[i] = v.(string)
	}

	new_array_bool := make([]bool, len(new_array_set))

	var del_array_ids []string
	var add_array_ids []string

	for i := 0; i < len(old_array_ids); i++ {
		exist := false
		for j := 0; j < len(new_array_ids); j++ {
			if new_array_ids[j] == old_array_ids[i] {
				exist = true
				new_array_bool[j] = true
				break
			}
		}
		if !exist {
			del_array_ids = append(del_array_ids, old_array_ids[i])
		}
	}
	for i := 0; i < len(new_array_bool); i++ {
		if !(new_array_bool[i]) {
			add_array_ids = append(add_array_ids, new_array_ids[i])
		}
	}

	return add_array_ids, del_array_ids
}

func substr(input string, start int, length int) string {
        asRunes := []rune(input)

        if start >= len(asRunes) {
                return ""
        }

        if start+length > len(asRunes) {
                length = len(asRunes) - start
        }

        return string(asRunes[start : start+length])
}

// To write/update the security object rotation_policy
func sobj_rotation_policy_write(rp map[string]interface{}) map[string]interface{} {
        rotation_policy := make(map[string]interface{})
        for k, v := range  rp{
            /* while sending the request, interval_days should be assigned as an integer.
               Hence it is converted to integer from the string.
            */
            if k == "interval_days" || k == "interval_months" {
                val, _ := strconv.Atoi(v.(string))
                rotation_policy[k] = val
            } else if k == "deactivate_rotated_key" {
                /* while sending the request, deactivate_rotated_key should be assigned as a boolean..
                   Hence it is converted to boolean from the string.
                */
                val, _ := strconv.ParseBool(v.(string))
                rotation_policy[k] = val
            } else {
                rotation_policy[k] = v
            }
        }
        return rotation_policy
}

// To read the security object rotation_policy
func sobj_rotation_policy_read(rp map[string]interface{}) map[string]interface{}  {
        rotation_policy := make(map[string]interface{})
        for k, v := range  rp{
            /* while reading the rotation_policy from terraform the interval_days attribute is assigned as float64 datatype.
               Hence it will be converted to string from float object.
            */
            if k == "interval_days" || k == "interval_months" {
                rotation_policy[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
            } else if k == "deactivate_rotated_key" {
                rotation_policy[k] = strconv.FormatBool(v.(bool))
            } else {
                rotation_policy[k] = v
            }
        }
        return rotation_policy
}