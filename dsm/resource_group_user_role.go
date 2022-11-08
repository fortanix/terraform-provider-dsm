package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Group
func resourceGroupUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGroupUserRole,
		ReadContext:   resourceReadGroupUserRole,
		UpdateContext: resourceUpdateGroupUserRole,
		DeleteContext: resourceDeleteGroupUserRole,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Group
func resourceCreateGroupUserRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := bind_group_user_role("add_groups", ctx, d, m)
	return diags
}

// [C]: Update Group
func resourceUpdateGroupUserRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := bind_group_user_role("mod_groups", ctx, d, m)
	return diags
}

// [R]: Read Group
func resourceReadGroupUserRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := dataSourceUserRead(ctx, d, m)
	return diags
}

// [D]: Delete Group
func resourceDeleteGroupUserRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := bind_group_user_role("del_groups", ctx, d, m)
	return diags
}

func bind_group_user_role(mode string, ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	group_name := d.Get("group_name").(string)
	user_email := d.Get("user_email").(string)
	role_name := d.Get("role_name").(string)
	role_id := role_name

	if err := d.Set("name", group_name); err != nil {
		return diag.FromErr(err)
	}
	dataSourceGroupRead(ctx, d, m)
	group_id := d.Get("group_id").(string)
	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("1 Group ID for group-user-role binding operation: %s", group_id))
	}

	if err := d.Set("user_email", user_email); err != nil {
		return diag.FromErr(err)
	}
	dataSourceUserRead(ctx, d, m)
	user_id := d.Get("user_id").(string)
	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("User ID for group-user-role binding operation: %s", user_id))
	}

	if role_name != "GROUPAUDITOR" && role_name != "GROUPADMINISTRATOR" {
		if err := d.Set("name", role_name); err != nil {
			return diag.FromErr(err)
		}
		diags := dataSourceRoleRead(ctx, d, m)
		if diags != nil {
			return diags
		}
		role_id = d.Get("role_id").(string)
	}

	main_object := make(map[string]interface{})
	sub_object := make(map[string]interface{})

	main_object["user_id"] = user_id
	sub_object[group_id] = []string{role_id}
	main_object[mode] = sub_object

	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("Role ID for group-user-role binding operation: %s", role_id))
		tflog.Warn(ctx, fmt.Sprintf("Main object for group-user-role binding operation: %s", main_object))
	}

	resp, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/users/%s", user_id), main_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}

	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[U]: API response for group-user-role binding operation: %s", resp_json))
	}

	d.SetId(resp["user_id"].(string))

	if err := d.Set("user_id", user_id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role_id", role_id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", group_id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
