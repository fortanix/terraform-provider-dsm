// **********
// Terraform Provider - SDKMS: resource: security object
// **********
//       - Author:    Ravi Gopal at fortanix dot com
//       - Version:   0.5.36
//       - Date:      24/07/2025
// **********

package dsm

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define User
func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateUser,
		ReadContext:   resourceReadUser,
		UpdateContext: resourceUpdateUser,
		DeleteContext: resourceDeleteUser,
		Description: "Creates a Fortanix DSM user. The returned resource object contains the UUID of the user for further references.",
		Schema: map[string]*schema.Schema{
			"user_email": {
				Description: "User's Email Id.",
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Description: "Description of a user.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Description: "User's Id.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Description: "Role of an account. Possible values are: \n" +
				 "   * `ACCOUNTADMINISTRATOR`\n" +
				 "   * `ACCOUNTMEMBER`\n" +
				 "   * `ACCOUNTAUDITOR` ",
				Type:     schema.TypeString,
				Required: true,
			},
			"account_role": {
				Description: "Role of an account. It also defines the state of a user.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"first_name": {
				Description: "User's first name",
				Type:     schema.TypeString,
				Optional: true,
				Default: nil,
			},
			"last_name": {
				Description: "User's last name",
				Type:     schema.TypeString,
				Optional: true,
				Default: nil,
			},
			"groups": {
				Description: "Add user to specific groups.\n\n" +
				"   * **Note:**  This parameter is essentially only for ACCOUNTMEMBER.\n",
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_verified": {
				Description: "If an email is verified.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_password": {
				Description: "If the user has password.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			"self_provisioned": {
				Description: "Is the user self self provisioned.",
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// Create
func resourceCreateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	user := map[string]interface{}{
		"user_email": d.Get("user_email").(string),
	}
	role := d.Get("role").(string)
	account_role := [1]string{role}
	user["account_role"] = account_role
	if first_name := d.Get("first_name").(string); len(first_name) > 0 {
		user["first_name"] = first_name
	}
	if last_name := d.Get("last_name").(string); len(last_name) > 0 {
		user["last_name"] = last_name
	}
	req, err := m.(*api_client).APICallBody("POST", dsm_endpoints["user_invite"], user)
	if err != nil {
		return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: POST %s: %v", dsm_endpoints["user_invite"], err), error_summary)
	}
	d.SetId(req["user_id"].(string))
	// Since creating a user and adding the groups to it can't be created in a single API,
	// So, patch call should be invoked to configure the groups
	if groups, ok := d.GetOk("groups"); ok && role == "ACCOUNTMEMBER" {
		patch_user := map[string]interface{}{
			"user_id": req["user_id"],
		}
		add_groups, err := unmarshalStringToJson(groups.(string))
		if err != nil {
			return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: %v", err), error_summary)
		}
		patch_user["add_groups"] = add_groups
		patch_url := fmt.Sprintf("%s/%s", dsm_endpoints["user"], req["user_id"].(string))
		_, err1 := m.(*api_client).APICallBody("PATCH", patch_url, patch_user)
		if err1 != nil {
			d.Set("groups", "")
			resourceDeleteUser(ctx, d, m)
			return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: PATCH %s: %v", patch_url, err1), error_summary)

		}
	}
	return resourceReadUser(ctx, d, m)
}

// Read
func resourceReadUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("%s/%s", dsm_endpoints["user"], d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: GET sys/v1/users/%s: %v", d.Id(), err), error_summary)
		}
		if err := d.Set("user_email", req["user_email"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("account_role", req["account_role"].([]interface{})); err != nil {
			return diag.FromErr(err)
		}
		if description, ok := req["description"]; ok {
			if err := d.Set("description", description.(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		// Only ACCOUNTMEMBER can add or delete groups.
		// Whereas both ACCOUNTADMINISTRATOR and ACCOUNTAUDITOR are part of all the groups.
		// So, setting the groups information in tf state for ACCOUNTADMINISTRATOR and ACCOUNTAUDITOR is not needed.
		// It also consumes memory.
		// Hence, groups information will be written only for an ACCOUNTMEMBER.
		if d.Get("role") == "ACCOUNTMEMBER" {
			if groups, ok := req["groups"]; ok {
				groups_string, err := json.Marshal(groups.(map[string]interface{}))
				if err != nil {
					if err := d.Set("groups", string(groups_string)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		if email_verified, ok := req["email_verified"]; ok {
			if err := d.Set("email_verified", email_verified.(bool)); err != nil {
				return diag.FromErr(err)
			}
		}
		if has_password, ok := req["has_password"]; ok {
			if err := d.Set("has_password", has_password.(bool)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

// Update
func resourceUpdateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	user := make(map[string]interface{})
	if d.HasChange("description") {
		user["description"] = d.Get("description").(string)
	}
	if d.HasChange("role") {
		account_role := [1]string{d.Get("role").(string)}
		user["account_role"] = account_role
	}
	if d.HasChange("groups") {
		old_groups, new_groups := d.GetChange("groups")
		old_groups_json := make(map[string][]string)
		old_groups_json, o_err := ConvertStringToJSONGeneric[map[string][]string](old_groups.(string))
		if o_err != nil {
		    return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: PATCH %s: %v", dsm_endpoints["user"], o_err), error_summary)
		}
		new_groups_json := make(map[string][]string)
		new_groups_json, n_err := ConvertStringToJSONGeneric[map[string][]string](new_groups.(string))
		if n_err != nil {
		    return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: PATCH %s: %v", dsm_endpoints["user"], n_err), error_summary)
		}
		var old_group_ids []interface{}
		for k, _ := range old_groups_json {
			old_group_ids = append(old_group_ids, k)
		}
		var new_group_ids []interface{}
		for k, _ := range new_groups_json {
			new_group_ids = append(new_group_ids, k)
		}
		// Get the groups to be added and deleted
		add_group_ids, del_group_ids := compute_add_and_del_arrays(old_group_ids, new_group_ids)
		if len(add_group_ids) > 0 {
			add_groups := make(map[string]interface{})
			for _, g_id := range add_group_ids {
				add_groups[g_id] = new_groups_json[g_id]
			}
			user["add_groups"] = add_groups
		}
		if len(del_group_ids) > 0 {
			del_groups := make(map[string]interface{})
			for _, g_id := range del_group_ids {
				del_groups[g_id] = old_groups_json[g_id]
			}
			user["del_groups"] = del_groups
		}
		// Get the groups to be modified
		mod_groups := make(map[string]interface{})
		for k, v1 := range old_groups_json {
			if v2, ok := new_groups_json[k]; ok {
				for i := range v1 {
					if v1[i] != v2[i] {
						mod_groups[k] = v2
					}
				}
			}
		}
		if len(mod_groups) > 0 {
			user["mod_groups"] = mod_groups
		}
	}
	patch_url := fmt.Sprintf("%s/%s", dsm_endpoints["user"], d.Id())
	_, err := m.(*api_client).APICallBody("PATCH", patch_url, user)
	if err != nil {
		return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: PATCH %s: %v", dsm_endpoints["user"], err), error_summary)
	}
	return resourceReadUser(ctx, d, m)
}

// Delete
func resourceDeleteUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	groups := d.Get("groups").(string)
	if len(groups) > 0 {
		user := make(map[string]interface{})
		del_groups, err := unmarshalStringToJson(groups)
		if err != nil {
			return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: %v", err), error_summary)
		}
		user["del_groups"] = del_groups
		patch_url := fmt.Sprintf("%s/%s", dsm_endpoints["user"], d.Id())
		_, p_err := m.(*api_client).APICallBody("PATCH", patch_url, user)
		if p_err != nil {
			return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: PATCH %s: %v", dsm_endpoints["user"], p_err), error_summary)
		}
	}
	_, _, err := m.(*api_client).APICall("DELETE", fmt.Sprintf(dsm_endpoints["user"] + "/%s" + "/accounts", d.Id()))
	if err != nil {
		return invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: DELETE %s: %v", dsm_endpoints["user"], err), error_summary)
	}
	d.SetId("")
	return nil
}
