package dsm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// [-] Define Security Object
func resourceSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateSobject,
		ReadContext:   resourceReadSobject,
		UpdateContext: resourceUpdateSobject,
		DeleteContext: resourceDeleteSobject,
		Description: "Creates a new security object. The returned resource object contains the UUID of the security object for further references.\n" +
		"A key value can be imported as a security object. This resource also can rotate or copy a security object.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The security object name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"dsm_name": {
			    Description: "The security object name.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
			    Description: "The security object group assignment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"obj_type": {
			    Description: "The security object type. Supported security objects:\n" +
			    "AES, DES, DES3, RSA, DSA, KCDSA, EC, ECKCDSA, ARIA, SEED and Tokenization(fpe).",
				Type:     schema.TypeString,
				Required: true,
			},
			"key_size": {
			    Description: "The security object size. It should not be given only when the obj_type is EC.\n" +
				"| obj_type | key_size | key_ops |\n" +
				"| -------- | -------- |-------- |\n" +
				"| `RSA` | 1024, 2048, 4096, 8192 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |\n" +
				"| `DSA` | 2048, 3072 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |\n" +
				"| `KCDSA` | 2048 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |\n" +
				"| `AES` | 128, 192, 256 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |\n" +
				"| `DES` | 56 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, APPMANAGEABLE, EXPORT |\n" +
				"| `DES3` | 112, 168 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |\n" +
				"| `ARIA` | 128, 192, 256 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |\n" +
				"| `SEED` | 128 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, EXPORT |\n",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"kid": {
			    Description: "The security object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
			    Description: "Account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},

			//"kcv": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"rotate": {
			    Description: "specify method to use for key rotation.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DSM", "ALL"}, true),
			},
			"rotate_from": {
			    Description: "Name of the security object to be rotated from.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator": {
			    Description: "The creator of the security object from Fortanix DSM.\n" +
			    "   * `user`: If the security object was created by a user, the computed value will be the matching user id.\n" +
			    "   * `app`: If the security object was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rotation_policy": {
				Description: "Policy to rotate a Security Object, configure the below parameters.\n" +
				"   * `interval_days`: Rotate the key for every given number of days.\n" +
				"   * `interval_months`: Rotate the key for every given number of months.\n" +
				"   * `effective_at`: Start of the rotation policy time.\n" +
				"   * `rotate_copied_keys`: Enable key rotation for copied keys.\n" +
				"   * **Note:** Either interval_days or interval_months should be given, but not both.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			// Unable to define links
			//"links": {
			//	Type:     schema.TypeMap,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeList,
			//	},
			//},
			"copied_to": {
			    Description: "List of security objects copied by the current security object.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_from": {
			    Description: "Security object that is copied to the current security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replacement": {
			    Description: "Replacement of a security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replaced": {
			    Description: "Replaced by a security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"pub_key": {
			    Description: "Public key (if ”RSA” obj_type is specified).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_pub_key": {
			    Description: "Open SSH public key (if ”RSA” obj_type is specified).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_metadata": {
			    Description: "The user defined security object attributes added to the key’s metadata from Fortanix DSM.",
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"fpe_radix": {
			    Description: "integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fpe": {
				Description: "FPE specific options. obj_type should be AES. It should be given in string format like below:\n" +
				"```This is a sample variable that specifies fpeOptions to create a Tokenization object that can tokenize credit card format data:\n" +
				"    variable " + "\"fpeOptionsExample\"" + " { \n" +
				"      type = any\n" +
				"      description = " + "\"The policy document. This is a JSON formatted string.\"" + "\n" +
				"      default = <<-EOF \n" +
				"              {\n" +
				"               " + "\"description\"" + ":" + "\"Credit card\"" + "\n" +
				"               " + "\"format\"" + ": {\n" +
				"               " + "\"char_set\"" + ": [\n" +
				"                    [\n" +
				"                     "+ "\"0\"" + ",\n" +
				"                     "+ "\"9\"" + "\n" +
				"                    ]\n" +
				"                  ],\n" +
				"                  " + "\"min_length\"" + ": 13,\n" +
				"                  " + "\"max_length\"" + ": 19,\n" +
				"                  " + "\"constraints\"" + ": {\n" +
				"                   " + "\"luhn_check\"" + ": true\n" +
				"                  }\n" +
				"              }\n" +
				"            }\n" +
				"            EOF\n" +
				"    }\n" +
				"\nThis is how we can reference this fpeOptions:\n" +
				"      fpe = var.fpeOptionsExample\n" +
				"\nRefer to the fpeOptions schema in https://www.fortanix.com/fortanix-restful-api-references/dsm for a better understanding of the fpe body.\n" +
				"```",
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_ops": {
			    Description: " The security object key permission from Fortanix DSM.\n" +
			    "   * Default is to allow all permissions except EXPORT",
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rsa": {
			    Description: "rsaOptions passed as a string (if ”RSA” obj_type is specified). The string should match the 'rsa' value in Post body while working with Fortanix Rest API. For Example:\n" +
			    "\n`rsa = " + "\"{\\" + "\"encryption_policy\\\"" + ":[{\\" + "\"padding\\\"" + ":{\\" + "\"RAW_DECRYPT\\\"" + ":{}}},{\\" + "\"padding\\\"" + ":{\\" + "\"OAEP\\\"" + ":{\\" + "\"mgf\\\"" + ":{\\" + "\"mgf1\\\"" + ":{\\" + "\"hash\\\"" + ":\\"+ "\"SHA1\\\""+ "}}}}}],\\"+ "\"signature_policy\\\"" + ":[{\\" + "\"padding\\\"" + ":{\\" + "\"PKCS1_V15\\\"" + ":{}}},{\\" + "\"padding\\\"" + ":{\\" + "\"PSS\\\"" + ":{\\" + "\"mgf\\\"" + ":{\\" + "\"mgf1\\\"" + ":{\\" + "\"hash\\\"" + ":\\" + "\"SHA384\\\"" + "}}}}}]}" + "\"" + "`",
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_key_justifications_policy": {
			    Description: "The security object key justification policies for GCP External Key Manager. The allowed permissions are:\n" +
			    "   * CUSTOMER_INITIATED_SUPPORT\n" +
			    "   * CUSTOMER_INITIATED_ACCESS\n" +
			    "   * GOOGLE_INITIATED_SERVICE\n" +
			    "   * GOOGLE_INITIATED_REVIEW\n" +
			    "   * GOOGLE_INITIATED_SYSTEM_OPERATION\n" +
			    "   * THIRD_PARTY_DATA_REQUEST\n" +
			    "   * REASON_NOT_EXPECTED\n" +
			    "   * REASON_UNSPECIFIED\n" +
			    "   * MODIFIED_CUSTOMER_INITIATED_ACCESS\n" +
			    "   * MODIFIED_GOOGLE_INITIATED_SYSTEM_OPERATION\n" +
			    "   * GOOGLE_RESPONSE_TO_PRODUCTION_ALERT\n",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"allowed_missing_justifications": {
			    Description: " Boolean value which allows missing justifications even if not provided to the security object. The values are True / False.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
			    Description: "The security object description.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
			    Description: "Whether the security object is enabled or disabled.\n" +
			    "   * The values are True/False.",
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"state": {
			    Description: "The state of the secret security object.\n" +
			    "   * Allowed states are: None, PreActive, Active, Deactivated, Compromised, Destroyed, Deleted.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_date": {
			    Description: "The security object expiry date in RFC format.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"elliptic_curve": {
				Description: "Standardized elliptic curve. It should be given only when the obj_type is EC or ECKCDSA.\n" +
				"| obj_type | Curve | key_ops |\n" +
				"| -------- | -------- |-------- |\n" +
				"| `EC` | SecP192K1, SecP224K1, SecP256K1  NistP192, NistP224, NistP256, NistP384, NistP521, X25519, Ed25519 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |\n" +
				"| `ECKCDSA` | SecP192K1, SecP224K1, SecP256K1  NistP192, NistP224, NistP256, NistP384, NistP521 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |\n",
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
			    Description: "Sobject content when importing content.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"subgroup_size": {
				Description: "Subgroup Size for DSA and ECKCDSA. The allowed Subgroup Sizes are 224 and 256.\n" +
				"| obj_type | subgroup_size | usage\n" +
				"| -------- | -------- | -------- |\n"+
				"| `DSA` | 224, 256| 224: When DSA key_size is 2048. 256: When DSA key_size is 2048 and 3072.\n" +
				"| `KCDSA` | 224, 256| 224, 256: When KCDSA key_size is 2048.\n",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hash_alg": {
			    Description: "Hashing Algorithm for KCDSA and ECKCDSA.\n" +
				"| obj_type | hash_alg |\n" +
				"| -------- | -------- |\n"+
				"| `ECKCDSA` | SHA1,SHA224, SHA256, SHA384, SHA521|\n" +
				"| `KCDSA` | SHA224, SHA256 |\n",
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [-]: Custom Functions
// contains: Need to validate whether a string exists in a []string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// this function takes a string in JSON format and unmarshals it. If the string is not in correct JSON format, it returns nil.
func unmarshalStringToJson(inputString string) (interface{}, error) {
	type mapFormat map[string]interface{}
	var inputMap mapFormat
	if err := json.Unmarshal([]byte(inputString), &inputMap); err != nil {
		return nil, err
	}

	return inputMap, nil
}

// createSO: Create Security Object
func createSO(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := "crypto/v1/keys"
	key_size := d.Get("key_size").(int)
	obj_type := d.Get("obj_type").(string)
	elliptic_curve := d.Get("elliptic_curve").(string)
	hash_alg := d.Get("hash_alg").(string)
	subgroup_size := d.Get("subgroup_size").(int)
	method := "POST"

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"obj_type":    obj_type,
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if _, ok := d.GetOk("value"); ok {
		security_object["value"] = d.Get("value").(string)
		method = "PUT"
	} else {
		if obj_type == "EC" || obj_type == "ECKCDSA" {
			if key_size > 0 || len(elliptic_curve) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Detail:   fmt.Sprintf("key_size should not be specified and elliptic_curve should be specified for %s", obj_type),
				})
				return diags
			} else {
				security_object["elliptic_curve"] = elliptic_curve
			}
		} else if key_size == 0 || len(elliptic_curve) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Detail:   fmt.Sprintf("key_size should be specified and elliptic_curve should not be specified for %s", obj_type),
			})
			return diags
		} else {
			security_object["key_size"] = key_size
		}
	}

	if rfcdate := d.Get("expiry_date").(string); len(rfcdate) > 0 {
		layoutRFC := "2006-01-02T15:04:05Z"
		layoutDSM := "20060102T150405Z"
		ddate, newerr := time.Parse(layoutRFC, rfcdate)
		if newerr != nil {
			return diag.FromErr(newerr)
		}
		security_object["deactivation_date"] = ddate.Format(layoutDSM)
	}

	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}
	if err := d.Get("rsa").(string); len(err) > 0 {
		rsa_obj, er := unmarshalStringToJson(d.Get("rsa").(string))
		if er != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid json string format for the field 'rsa'.",
				Detail:   fmt.Sprintf("[E]: Input: rsa: %s", err),
			})
			return diags
		}
		security_object["rsa"] = rsa_obj
	}
	allowed_key_justifications_policy, ok := d.GetOk("allowed_key_justifications_policy")
	allowed_missing_justifications, ok2 := d.GetOkExists("allowed_missing_justifications")

	if ok && ok2 {
		if allowed_key_justifications_policy != nil && allowed_missing_justifications != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow":                allowed_key_justifications_policy,
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	} else if ok {
		if allowed_key_justifications_policy != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow": allowed_key_justifications_policy,
			}
		}
	} else if ok2 {
		if allowed_missing_justifications != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	}

	// Ensuring that only one of these options (`fpe`, `fpe_radix`) is specified in the Terraform configuration to maintain backward compatibility.
	// This prevents issues for existing users of fpe_radix.
	// This logic was added in v0.5.30 to support the transition from `fpe_radix` to `fpe` for new users while maintaining support for existing configurations.
	if d.Get("fpe").(string) != "" && d.Get("fpe_radix").(int) != 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "only one of these two can be given in the Terraform configuration: fpe, fpe_radix. This check ensures backward compatibility for users previously using 'fpe_radix'. New users are encouraged to use the 'fpe' object.",
		})
		return diags
	}
	if fpe_policy := d.Get("fpe").(string); len(fpe_policy) > 0 {
		security_object["fpe"] = json.RawMessage(d.Get("fpe").(string))
	}
	if err := d.Get("fpe_radix"); err != 0 {
		security_object["fpe"] = map[string]interface{}{
			"radix": d.Get("fpe_radix").(int),
		}
	}
	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = err
	}
	if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
	}

	if len(hash_alg) > 0 && obj_type == "KCDSA" {
		kcdsa := make(map[string]interface{})
		kcdsa["hash_alg"] = hash_alg
		kcdsa["subgroup_size"] = subgroup_size
		security_object["kcdsa"] = kcdsa
	} else if len(hash_alg) > 0 && obj_type == "ECKCDSA" {
		eckcdsa := make(map[string]interface{})
		eckcdsa["hash_alg"] = hash_alg
		security_object["eckcdsa"] = eckcdsa
	} else if obj_type == "DSA" {
		dsa := make(map[string]interface{})
		dsa["subgroup_size"] = subgroup_size
		security_object["dsa"] = dsa
	}

	if err := d.Get("rotate").(string); len(err) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
		endpoint = "crypto/v1/keys/rekey"
	}
	req, err := m.(*api_client).APICallBody(method, endpoint, security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST %s: %v", endpoint, err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return diags
}

// [C]: Terraform Func: resourceCreateSobject
func resourceCreateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		if rotate_from := d.Get("rotate_from").(string); len(rotate_from) <= 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   "[E]: API: GET crypto/v1/keys/rekey: 'rotate_from' missing",
			})
			return diags
		}
	}

	if err := createSO(ctx, d, m); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err),
		})
		return diags
	}

	return resourceReadSobject(ctx, d, m)
}

// [R]: Terraform Func: resourceReadSobject
func resourceReadSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err),
			})
			return diags
		}

		if err := d.Set("dsm_name", req["name"].(string)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("obj_type", req["obj_type"].(string)); err != nil {
			return diag.FromErr(err)
		}
		obj_type := req["obj_type"].(string)
		if req["origin"] != "External" {
			if _, ok := req["key_size"]; ok {
				if err := d.Set("key_size", int(req["key_size"].(float64))); err != nil {
					return diag.FromErr(err)
				}
			}
			if _, ok := req["elliptic_curve"]; ok {
				if err := d.Set("elliptic_curve", req["elliptic_curve"].(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if _, ok := req["google_access_reason_policy"]; ok {
			google_access_reason_policy := req["google_access_reason_policy"].(map[string]interface{})
			if err := d.Set("allowed_key_justifications_policy", google_access_reason_policy["allow"]); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("allowed_missing_justifications", google_access_reason_policy["allow_missing_reason"]); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("kid", req["kid"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["pub_key"]; ok {
			if err := d.Set("pub_key", req["pub_key"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", req["creator"]); err != nil {
			return diag.FromErr(err)
		}
		//if err := d.Set("links", req["links"]); err != nil {
		//	return diag.FromErr(err)
		//}
		// FYOO: Fix this later - some wierd reaction to TypeList/TypeMap within TF
		if err := d.Set("copied_to", req["copied_to"]); err != nil {
			return diag.FromErr(err)
		}

		if _, ok := req["links"]; ok {
			if links := req["links"].(map[string]interface{}); len(links) > 0 {
				if _, copiedToExists := req["links"].(map[string]interface{})["copiedTo"]; copiedToExists {
					if err := d.Set("copied_to", req["links"].(map[string]interface{})["copiedTo"].([]interface{})); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, copiedFromExists := req["links"].(map[string]interface{})["copiedFrom"]; copiedFromExists {
					if err := d.Set("copied_from", req["links"].(map[string]interface{})["copiedFrom"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacementExists := req["links"].(map[string]interface{})["replacement"]; replacementExists {
					if err := d.Set("replacement", req["links"].(map[string]interface{})["replacement"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacedExists := req["links"].(map[string]interface{})["replaced"]; replacedExists {
					if err := d.Set("replaced", req["links"].(map[string]interface{})["replaced"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
			return diag.FromErr(err)
		}
		if err := req["fpe"]; err != nil {
			if req["fpe"].(map[string]interface{})["radix"] != nil && d.Get("fpe_radix") != nil {
				if err := d.Set("fpe_radix", int(req["fpe"].(map[string]interface{})["radix"].(float64))); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err := d.Set("fpe", (d.Get("fpe").(string)) ); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		// FYOO: Fix TypeList sorting error
		key_ops := make([]string, len(req["key_ops"].([]interface{})))
		if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
			if len(d.Get("key_ops").([]interface{})) == len(req["key_ops"].([]interface{})) {
				for idx, key_op := range d.Get("key_ops").([]interface{}) {
					key_ops[idx] = fmt.Sprint(key_op)
				}
			} else {
				req_key_ops := make([]string, len(req["key_ops"].([]interface{})))
				for idx, key_op := range req["key_ops"].([]interface{}) {
					req_key_ops[idx] = fmt.Sprint(key_op)
				}
				final_idx := 0
				for _, key_op := range d.Get("key_ops").([]interface{}) {
					if contains(req_key_ops, fmt.Sprint(key_op)) {
						key_ops[final_idx] = fmt.Sprint(key_op)
						final_idx = final_idx + 1
					}
				}
			}
		} else {
			for idx, key_op := range req["key_ops"].([]interface{}) {
				key_ops[idx] = fmt.Sprint(key_op)
			}
		}
		if err := d.Set("key_ops", key_ops); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["description"]; ok {
			if err := d.Set("description", req["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("state", req["state"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if rfcdate, ok := req["deactivation_date"]; ok {
			// FYOO: once it's set, you can't remove deactivation date
			layoutRFC := "2006-01-02T15:04:05Z"
			layoutDSM := "20060102T150405Z"
			ddate, newerr := time.Parse(layoutDSM, rfcdate.(string))
			if newerr != nil {
				return diag.FromErr(newerr)
			}
			if newerr = d.Set("expiry_date", ddate.Format(layoutRFC)); newerr != nil {
				return diag.FromErr(newerr)
			}
		}
		if err := req["obj_type"].(string); err == "RSA" {
			openssh_pub_key, err := PublicPEMtoOpenSSH([]byte(req["pub_key"].(string)))
			if err != nil {
				return err
			} else {
				if err := d.Set("ssh_pub_key", openssh_pub_key); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if obj_type == "DSA" {
			if _,ok := req["dsa"]; ok {
				dsa := req["dsa"].(map[string]interface{})
				if err := d.Set("subgroup_size", dsa["subgroup_size"]); err != nil {
					return diag.FromErr(err)
				}
			}
		} else if obj_type == "KCDSA" {
			if _,ok := req["kcdsa"]; ok {
				kcdsa := req["kcdsa"].(map[string]interface{})
				if err := d.Set("subgroup_size", kcdsa["subgroup_size"]); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("hash_alg", kcdsa["hash_alg"]); err != nil {
					return diag.FromErr(err)
				}
			}
		} else if obj_type == "ECKCDSA" {
			if _,ok := req["eckcdsa"]; ok {
				eckcdsa := req["eckcdsa"].(map[string]interface{})
				if err := d.Set("hash_alg", eckcdsa["hash_alg"]); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		}
		if _, ok := req["rotation_policy"]; ok {
			rotation_policy := sobj_rotation_policy_read(req["rotation_policy"].(map[string]interface{}))
			if err := d.Set("rotation_policy", rotation_policy); err != nil {
				return diag.FromErr(err)
			}
		}

		// FYOO: clear values that are irrelevant
		d.Set("rotate", "")
		d.Set("rotate_from", "")
		return diags
}

// [U]: Terraform Func: resourceUpdateSobject
func resourceUpdateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var has_changed = false

	// already has been replaced so "rotate" and "rotate_from" does not apply
	_, replacement := d.GetOk("replacement")
	_, replaced := d.GetOk("replaced")
	if replacement || replaced {
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	var security_object = map[string]interface{}{
		"kid": d.Get("kid").(string),
	}
	if d.HasChange("rsa") {
		rsa_obj, err := unmarshalStringToJson(d.Get("rsa").(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid json string format for the field 'rsa'.",
				Detail:   fmt.Sprintf("[E]: Input: rsa: %s", err),
			})
			return diags
		}
		security_object["rsa"] = rsa_obj
		has_changed = true
	}
	if d.HasChanges("allowed_key_justifications_policy", "allowed_missing_justifications") {

		google_access_reason_policy := make(map[string]interface{})

		google_access_reason_policy["allow"] = d.Get("allowed_key_justifications_policy")
		google_access_reason_policy["allow_missing_reason"] = d.Get("allowed_missing_justifications")

		has_changed = true

		security_object["google_access_reason_policy"] = google_access_reason_policy
	}
	if d.HasChange("key_ops") {
		security_object["key_ops"] = d.Get("key_ops")
		has_changed = true
	}
	if d.HasChange("description") {
		security_object["description"] = d.Get("description")
		has_changed = true
	}
	if d.HasChange("name") {
		security_object["name"] = d.Get("name")
		has_changed = true
	}
	if d.HasChange("custom_metadata") {
		security_object["custom_metadata"] = d.Get("custom_metadata").(map[string]interface{})
		has_changed = true
	}
	if d.HasChange("rotation_policy") {
		rotation_policy := d.Get("rotation_policy").(map[string]interface{})
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
		has_changed = true
	}
	if d.HasChange("fpe") {
		old_fpe, new_fpe := d.GetChange("fpe")
		d.Set("fpe", old_fpe)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "fpe cannot modify on update",
			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: fpe cannot change on update. Please retain it to old value: %s -> %s", old_fpe, new_fpe),
		})
		return diags
	}
	if d.HasChange("hash_alg") {
		old_ha, new_ha := d.GetChange("hash_alg")
		d.Set("hash_alg", old_ha)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "hash_alg cannot modify on update",
			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: hash_alg cannot change on update. Please retain it to old value: %s -> %s", old_ha, new_ha),
		})
		return diags
	}
	if d.HasChange("subgroup_size") {
		old_sz, new_sz := d.GetChange("dsa")
		d.Set("subgroup_size", old_sz)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "dsa cannot modify on update",
			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: subgroup_size cannot change on update. Please retain it to old value: %s -> %s", old_sz, new_sz),
		})
		return diags
	}
	if d.HasChange("elliptic_curve") {
		old_ec, new_ec := d.GetChange("elliptic_curve")
		d.Set("elliptic_curve", old_ec)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "elliptic_curve cannot modify on update",
			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: elliptic_curve cannot change on update. Please retain it to old value: %s -> %s", old_ec, new_ec),
		})
		return diags
	}

	if has_changed {
		if debug_output {
			tflog.Warn(ctx, "Sobject has changed, calling API")
		}
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), security_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			return diags
		}

		key_ops := make([]string, len(req["key_ops"].([]interface{})))
		if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
			if len(d.Get("key_ops").([]interface{})) == len(req["key_ops"].([]interface{})) {
				for idx, key_op := range d.Get("key_ops").([]interface{}) {
					key_ops[idx] = fmt.Sprint(key_op)
				}
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Error in processing DSM provider API response key_ops",
					Detail:   "[E]: API: PATCH crypto/v1/keys: Sync issue from State and DSM",
				})
				return diags
			}
		} else {
			for idx, key_op := range req["key_ops"].([]interface{}) {
				key_ops[idx] = fmt.Sprint(key_op)
			}
		}
		if err := d.Set("key_ops", key_ops); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceReadSobject(ctx, d, m)
}

// [D]: Terraform Func: resourceDeleteSobject
func resourceDeleteSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))

	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %v", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
