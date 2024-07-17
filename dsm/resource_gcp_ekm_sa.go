// **********
// Terraform Provider - DSM: resource: gcp_ekm_sa
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      27/07/2021
//       - Changelog:
//                  - Initial release to support resource for app with GCP EKM SA
// **********

package dsm

import (
	"context"
	"fmt"

	//	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceGcpEkmSa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGcpEkmSa,
		ReadContext:   resourceReadGcpEkmSa,
		UpdateContext: resourceUpdateGcpEkmSa,
		DeleteContext: resourceDeleteGcpEkmSa,
		Description: "Returns the Fortanix DSM Google EKM app from the cluster as a resource.",
		Schema: map[string]*schema.Schema{
			// service account name = app name
			"name": {
				Description: "The Google service account name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Description: "The unique ID of the app from Terraform.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
				Description: "The Fortanix DSM group ID to be mapped to the app by default.",
				Type:     schema.TypeString,
				Required: true,
			},
			"acct_id": {
				Description: "The account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the security object from Fortanix DSM.\n" +
				"   * `user`: If the security object was created by a user, the computed value will be the matching user ID.\n" +
				"   * `app`: If the security object was created by an app, the computed value will be the matching app ID.",
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
			"other_group": {
				Description: "The Fortanix DSM group object ID the app needs to be assigned to. If you want to delete the existing groups from an app, remove the IDs during update.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"other_group_permissions": {
				Description: "If you want to change the default permissions of a new group.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"mod_group_permissions": {
				Description: "To modify the permissions of any existing group.",
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

// [C]: Create GcpEkmSa
func resourceCreateGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		"default_group": d.Get("default_group").(string),
		"app_type":    "default",
		"description": d.Get("description").(string),
	}
	// add groups and it's permissions
	formAddGroups(d, app_object)
	// add GCP specifics
	gcp_perm := map[string]interface{}{
		"allow":                []string{"CUSTOMER_INITIATED_SUPPORT", "CUSTOMER_INITIATED_ACCESS", "GOOGLE_INITIATED_SERVICE", "GOOGLE_INITIATED_REVIEW", "GOOGLE_INITIATED_SYSTEM_OPERATION", "THIRD_PARTY_DATA_REQUEST", "REASON_UNSPECIFIED", "REASON_NOT_EXPECTED", "MODIFIED_CUSTOMER_INITIATED_ACCESS"},
		"allow_missing_reason": true,
	}

	app_object["credential"] = map[string]interface{}{
		"googleserviceaccount": gcp_perm,
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
	return resourceReadGcpEkmSa(ctx, d, m)
}

// [R]: Read GcpEkmSa
func resourceReadGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
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
	}
	return diags
}

// [U]: Update GcpEkmSa
func resourceUpdateGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := make(map[string]interface{})
	app_object["description"] = d.Get("description")
	app_object["name"] = d.Get("name")
	if d.HasChange("default_group") {
		if default_group := d.Get("default_group").(string); len(default_group) > 0 {
			app_object["default_group"] = d.Get("default_group")
		}
	}
	//Modifies the existing groups
	if d.HasChange("mod_group_permissions") {
		err := getChangesInGroupPermissions(d, app_object)
		if err != nil {
			return err
		}
	}
	if d.HasChange("other_group") {
		getChangesInOtherGroups(d, app_object)
	}
	if d.HasChange("description") {
		app_object["description"] = d.Get("description")
	}
	if len(app_object) > 0 {
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/apps/%s", d.Id()), app_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v", d.Get("mod_group_permissions")),
			})
			return diags
		}
		d.SetId(req["app_id"].(string))
	}
	return resourceReadGcpEkmSa(ctx, d, m)
}

// [D]: Delete GcpEkmSa
func resourceDeleteGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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