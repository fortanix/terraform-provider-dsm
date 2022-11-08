// **********
// Terraform Provider - DSM: resource: app
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.3
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateApp,
		ReadContext:   resourceReadApp,
		UpdateContext: resourceUpdateApp,
		DeleteContext: resourceDeleteApp,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"other_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acct_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"new_credential": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"other_group_permissions": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"mod_group_permissions": {
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
func resourceCreateApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		"default_group": d.Get("default_group").(string),
		//"add_groups": map[string]interface{}{
		//	d.Get("default_group").(string): []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"},
		//},
		"app_Type":    "default",
		"description": d.Get("description").(string),
	}

	add_group_perms := form_group_permissions(d.Get("other_group_permissions"))
	app_add_group := make(map[string]interface{})
	if err := d.Get("other_group").([]interface{}); len(err) > 0 {
		for _, group_id := range d.Get("other_group").([]interface{}) {
			if perms, ok := add_group_perms[group_id.(string)]; ok {
				app_add_group[group_id.(string)] = perms
			} else {
				app_add_group[group_id.(string)] = []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"}
			}
		}
	}

	if perms, ok := add_group_perms[d.Get("default_group").(string)]; ok {
		app_add_group[d.Get("default_group").(string)] = perms
	} else {
		app_add_group[d.Get("default_group").(string)] = []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"}
	}

	app_object["add_groups"] = app_add_group

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
	return resourceReadApp(ctx, d, m)
}

// [R]: Read App
func resourceReadApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if err := d.Set("credential", base64.StdEncoding.EncodeToString([]byte(d.Id()+":"+req["credential"].(map[string]interface{})["secret"].(string)))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("new_credential", false); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// [U]: Update App
func resourceUpdateApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if d.Get("new_credential").(bool) {
		reset_secret := map[string]interface{}{
			"credential_migration_period": nil,
		}

		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/apps/%s/reset_secret", d.Id()), reset_secret)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
			})
			return diags
		}
	}
	//Modified by Ravi Gopal
	app_object := make(map[string]interface{})
	app_object["description"] = d.Get("description")

	if d.HasChange("default_group") {
		if default_group := d.Get("default_group").(string); len(default_group) > 0 {
			app_object["default_group"] = d.Get("default_group")
		}
	}
	if d.HasChange("other_group") {
		old_group, new_group := d.GetChange("other_group")
		// compute_add_and_del_arrays function is in common.go
		add_group_ids, del_group_ids := compute_add_and_del_arrays(old_group, new_group)
		//Add the groups to be deleted
		if len(del_group_ids) > 0 {
			app_object["del_groups"] = del_group_ids
		}
		//Add the new groups
		if len(add_group_ids) > 0 {
			add_group_perms := form_group_permissions(d.Get("other_group_permissions"))
			app_add_group := make(map[string]interface{})
			for i := 0; i < len(add_group_ids); i++ {
				if perms, ok := add_group_perms[add_group_ids[i]]; ok {
					app_add_group[add_group_ids[i]] = perms
				} else {
					app_add_group[add_group_ids[i]] = []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"}
				}
			}
			app_object["add_groups"] = app_add_group
		}
	}
	if d.HasChange("description") {
		app_object["description"] = d.Get("description")
	}
	//Modifies the existing groups
	if d.HasChange("mod_group_permissions") {
		if mod_group := d.Get("mod_group_permissions").(map[string]interface{}); len(mod_group) > 0 {
			app_mod_group := make(map[string]interface{})
			//if default_group has changes in permissions
			default_group := d.Get("default_group").(string)
			if perms, ok := mod_group[default_group]; ok {
				app_mod_group[default_group] = strings.Split(perms.(string), ",")
				delete(mod_group, default_group)
			}
			//checking whether all the group_ids from mod_group_permissions exists in other groups or not
			//if not it will ignore the mod_group_permissions of the unavailable group_id
			var other_group_latest []string
			if err := d.Get("other_group").([]interface{}); len(err) > 0 {
				for _, group_id := range d.Get("other_group").([]interface{}) {
					other_group_latest = append(other_group_latest, group_id.(string))
				}
			}
			for i := 0; i < len(other_group_latest); i++ {
				if perms, ok := mod_group[other_group_latest[i]]; ok {
					app_mod_group[other_group_latest[i]] = strings.Split(perms.(string), ",")
					delete(mod_group, other_group_latest[i])
				}
			}
			if len(mod_group) > 0 {
				var unavailable_group_ids []string
				for group_id := range mod_group {
					unavailable_group_ids = append(unavailable_group_ids, group_id)
				}
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "All the group_ids available in mod_group are not available in other_group.Please correct them.",
					Detail:   fmt.Sprintf("[E]: Input: mod_group: Please remove the group_ids from mod_group those are not part of other_group. \n Following group_ids are not available in other_group:\n %v", unavailable_group_ids),
				})
				if old_group, new_group := d.GetChange("other_group"); len(new_group.([]interface{})) > 0 {
					d.Set("other_group", old_group)
				}
				return diags
			}
			app_object["mod_groups"] = app_mod_group
		}
	}

	if len(app_object) > 0 {
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/apps/%s", d.Id()), app_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v", err),
			})

			return diags
		}
		d.SetId(req["app_id"].(string))
	}
	return resourceReadApp(ctx, d, m)
}

// [D]: Delete App
func resourceDeleteApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

// form the group permissions - Ravi Gopal
func form_group_permissions(permissions interface{}) map[string]interface{} {
	add_group_perms := make(map[string]interface{})
	if group_perms := permissions.(map[string]interface{}); len(group_perms) > 0 {
		for group_id, permissions := range group_perms {
			permissions_list := strings.Split(permissions.(string), ",")
			add_group_perms[group_id] = permissions_list
		}
	}

	return add_group_perms
}