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
	"time"

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

// To read the security object rotation_policy
func set_lms_read_sobject(lms map[string]interface{}) map[string]interface{}  {
        lms_data := make(map[string]interface{})
        for k, v := range  lms{
            /* while reading the lms value from terraform the interval_days attribute is assigned as float64 datatype.
               Hence it will be converted to string from float object.
            */
            if k == "node_size" || k == "l1_height" || k == "l2_height" {
                lms_data[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
             }
        }
        return lms_data
}


// Compare two array strings whether they are equal or not irrespective of the order.
func compTwoArrays(x interface{}, y interface{}) bool{
	x_array_set := x.([]interface{})
	y_array_set := y.([]interface{})

	xMap := make(map[string]int)
	yMap := make(map[string]int)
	for _, xElem := range x_array_set {
		xMap[xElem.(string)]++
	}
	for _, yElem := range y_array_set {
		yMap[yElem.(string)]++
	}
	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}

// return diagnostics without summary
func invokeErrorDiagsNoSummary(detail string) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Detail:   detail,
	})
	return diags
}

// return diagnostics with summary
func invokeErrorDiagsWithSummary(detail string, summary string) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary: summary,
		Detail:   detail,
	})
	return diags
}

// undo the tf state value from new to old if there is a failure
func undoTFstate(param_type string, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	old_val, new_val := d.GetChange(param_type)
	d.Set(param_type, old_val)
	diags = append(diags, diag.Diagnostic{
		  Severity: diag.Error,
		  Summary:  param_type + " cannot modify on update",
		  Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s cannot change on update. Please retain it to old value: %s -> %s", param_type, old_val, new_val),
	})
	return diags
}

// expiry_date create and update
func parseTimeToDSM(expiry_date string) (string, diag.Diagnostics){
	layoutRFC := "2006-01-02T15:04:05Z"
	layoutDSM := "20060102T150405Z"
	ddate, newerr := time.Parse(layoutRFC, expiry_date)
	if newerr != nil {
		return "", diag.FromErr(newerr)
	}
	return ddate.Format(layoutDSM), nil
}

/* Set the key_ops in tf state.
This function is required, request key_ops order and response key_ops order mighet differ.
Hence, during `terraform plan`, if there are no changes in key_ops, it should not show any changes.
*/
func setKeyOpsTfState(d *schema.ResourceData, key_ops interface{}) diag.Diagnostics{
	is_same_key_ops := compTwoArrays(key_ops, d.Get("key_ops"))
	if is_same_key_ops {
		if err := d.Set("key_ops", d.Get("key_ops")); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("key_ops", key_ops); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

/*
A BYOK security object can be deleted only when it is in Destroyed state.
*/
func deleteBYOKDestroyedSobject(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	/*
	BYOK object can be deleted only after destruction.
	*/
	error_summary := "[DSM SDK] Unable to call DSM provider API client"
	if d.Get("state").(string) == "Destroyed" {
		_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
		if (err != nil) && (statuscode != 204) {
		    return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %v", err))
		}
		d.SetId("")
	} else {
		return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: cannot be deleted as security object is in the state %s", d.Get("state").(string)))
	}
	return nil
}

/*
Throw a warning, if schedule_deletion fails after create.
*/
func scheduleDeletionWarning(d *schema.ResourceData, err diag.Diagnostics) diag.Diagnostics {
    var diags diag.Diagnostics
    error_summary := fmt.Sprintf("[DSM SDK] Creation of a security object is successful, but failed to schedule the deletion of a security object %s", d.Get("name").(string))
    d.Set("schedule_deletion", false)
    diags = append(diags, diag.Diagnostic{
        Severity: diag.Warning,
        Summary:  error_summary,
        Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/%s/schedule_deletion, %v", error_summary, d.Id(), err),
    })
    return diags
}