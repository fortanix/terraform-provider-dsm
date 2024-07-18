// **********
// Terraform Provider - DSM: resource: app
// **********
//       - Author:    ravigopal at fortanix dot com
//       - Version:   0.5.28
//       - Date:      27/12/2023
// **********

package dsm

import (
	"context"
	"fmt"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
This resource supports non-api key appications.
As of now, this supports the following:
1. AWS XKS
2. AWS IAM
3. Certificate
4. Trusted CA(dns_name and ip)

In future we can add other applications as well

*/

// [-] Define App
func resourceAppNonAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAppNonAPIKey,
		ReadContext:   resourceReadAppNonAPIKey,
		UpdateContext: resourceUpdateAppNonAPIKey,
		DeleteContext: resourceDeleteAppNonAPIKey,
		Description: "Creates a non API key app. The returned resource object contains the UUID of the app for further references. Default permissions of any group can be modified. Using dsm_app_non_api_key following " +
		"apps can be created:\n" +
		"   * `awsxks`\n" +
		"   * `awsiam`\n" +
		"   * `certificate`\n" +
		"   * `trustedca`\n",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The Fortanix DSM App name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
			    Description: "The unique ID of the app.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
			    Description: "The Fortanix DSM group object id to be mapped to the app by default.",
				Type:     schema.TypeString,
				Required: true,
			},
			"other_group": {
			    Description: "The Fortanix DSM group object id the app needs to be assigned to. If you want to delete the existing groups from an app, remove the ids during update.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acct_id": {
			    Description: "The account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the app from Fortanix DSM.\n" +
				"   * `user`: If the app was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the app was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
			    Description: "The description of the app.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"authentication_method": {
			    Description: "The Authentication type of an app.\n" +
				"   * `type`:  Following authentication types are supported.\n" +
				"       * awsxks, awsiam, certificate and trustedca.\n" +
				"   * `certificate`: Certificate value, this should be configured when the type is certificate.\n" +
				"   * `ca_certificate`: CA certificate value, this should be configured when the type is trustedca.\n" +
				"   One of the following parameters should be given when the type is trustedca.\n" +
				"   * `ip_address`:  IP address value for trusted ca.\n" +
				"   * `dns_name`:  DNS name for trusted ca.\n" +
				"   **Note**: For more details refer the above examples.",
				Type:      schema.TypeMap,
				Required:  true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"credential": {
			    Description: "The Fortanix DSM App credentials. When the authentication method is awsxks, " +
			    "AWSXKS access and secret keys will be mapped here.",
				Type:      schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"other_group_permissions": {
			    Description: "Incase if you want to change the default permissions of a new group that includes default group. Please refer the example.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"mod_group_permissions": {
			    Description: "To modify the permissions of any existing group that includes default group. Please refer the example.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create App
func resourceCreateAppNonAPIKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		"default_group": d.Get("default_group").(string),
		"app_Type":    "default",
		"description": d.Get("description").(string),
	}
	// add groups and it's permissions
	formAddGroups(d, app_object)
	if am := d.Get("authentication_method").(map[string]interface{}); len(am) > 0 {
		formCredential(d, app_object, am)
	}
	req, err := m.(*api_client).APICallBody("POST", "sys/v1/apps", app_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v", err),
		})
		return diags
	}

	d.SetId(req["app_id"].(string))
	return resourceReadAppNonAPIKey(ctx, d, m)
}

// [R]: Read App
func resourceReadAppNonAPIKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, _, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps: %v", err),
		})
		return diags
	}
	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_id", req["app_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_group", req["default_group"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creator", req["creator"]); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := req["description"]; ok {
		if err := d.Set("description", req["description"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	req, _, err = m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s/credential", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
		})
		return diags
	}
	credential := make(map[string]interface{})
	req_credential := req["credential"].(map[string]interface{})
	if val, ok := req_credential["awsxks"].(map[string]interface{}); ok {
		credential["access_key_id"] =  base64.StdEncoding.EncodeToString([]byte(val["access_key_id"].(string)))
		credential["secret_key"] = base64.StdEncoding.EncodeToString([]byte(val["secret_key"].(string)))
		credential["path_prefix"] = "/crypto/v1/apps/" + d.Id() + "/aws"
	} else if val, ok := req_credential["certificate"].(map[string]interface{}); ok {
		credential["certificate"] = val["certificate"].(string)
	}
	if err := d.Set("credential", credential); err != nil {
			return diag.FromErr(err)
	}

	return diags
}

// [U]: Update App
func resourceUpdateAppNonAPIKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := make(map[string]interface{})
	app_object["description"] = d.Get("description")
	app_object["name"] = d.Get("name")
	if d.HasChange("default_group") {
		if default_group := d.Get("default_group").(string); len(default_group) > 0 {
			app_object["default_group"] = d.Get("default_group")
		}
	}
	if d.HasChange("other_group") {
		getChangesInOtherGroups(d, app_object)
	}
	if d.HasChange("description") {
		app_object["description"] = d.Get("description")
	}
	//Modifies the existing groups
	if d.HasChange("mod_group_permissions") {
		err := getChangesInGroupPermissions(d, app_object)
		if err != nil {
			return err
		}
	}
	if d.HasChange("authentication_method") {
		am := d.Get("authentication_method").(map[string]interface{})
		formCredential(d, app_object, am)
	}
	if len(app_object) > 0 {
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/apps/%s", d.Id()), app_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v -%v", err, app_object),
			})

			return diags
		}
		d.SetId(req["app_id"].(string))
	}
	return resourceReadAppNonAPIKey(ctx, d, m)
}

// [D]: Delete App
func resourceDeleteAppNonAPIKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/apps: %v", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}

// authentication method
func formCredential(d *schema.ResourceData, app_object map[string]interface{}, am map[string]interface{}) {
	authentication_method := make(map[string]interface{})
	credential := make(map[string]interface{})
	credential_type := ""
	for k, v := range  am{
		if k == "type" {
			if v == "awsxks" {
				credential_type = "awsxks"
			} else if v == "awsiam" {
				credential_type = "awsiam"
			} else if v == "certificate" {
				credential_type = "certificate"
			} else if v == "trustedca" {
				credential_type = "trustedca"
			}
		} else if k == "certificate" {
			authentication_method["certificate"] = v.(string)
		} else if k == "ca_certificate" {
			credential["ca_certificate"] = v.(string)
		} else if k == "dns_name" {
			dns_name := make(map[string]string)
			dns_name["dns_name"] = v.(string)
			credential["subject_general"] = dns_name
		} else if k == "ip_address" {
			ip_address := make(map[string]string)
			ip_address["ip_address"] = v.(string)
			credential["subject_general"] = ip_address
		}
	}
	if credential_type != "certificate" {
		authentication_method[credential_type] = credential
	}
	app_object["credential"] = authentication_method
}