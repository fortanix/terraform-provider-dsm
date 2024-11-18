// **********
// Terraform Provider - DSM: resource: app
// **********
//       - Author:    ravigopal at fortanix dot com
//       - Version:   0.5.34
//       - Date:      15/11/2024
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceAppAssignGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAppAssignGroups,
		ReadContext:   resourceReadAppAssignGroups,
		UpdateContext: resourceUpdateAppAssignGroups,
		DeleteContext: resourceDeleteAppAssignGroups,
		Description: "Assigns new DSM group(s) to an existing DSM App. It assigns the default permissions to the newly added groups.\n" +
		"This resource is to only assign the groups to an App. Removing the groups from an App or modifying any other parameter is not possible.",
		Schema: map[string]*schema.Schema{
			"app_name": {
				Description: "The Fortanix DSM App name.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Description: "The unique ID of the app.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groups": {
				Description: "List of new DSM group IDs to be assigned to an App.",
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: PATCH App to add new groups
func resourceCreateAppAssignGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpoint := "sys/v1/apps"
	app_id := ""
	// Gets the App details through the app name or app id
	if app_name := d.Get("app_name").(string); len(app_name) > 0 {
		req, err := m.(*api_client).APICallList("GET", endpoint)
		if err != nil {
			return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: GET %s: %v", endpoint, err))
		}
		for _, data := range req {
			if data.(map[string]interface{})["name"].(string) == app_name {
				app_id = data.(map[string]interface{})["app_id"].(string)
				break
			}
		}
		endpoint += "/" + app_id
	} else if tf_app_id := d.Get("app_id").(string); len(tf_app_id) > 0 {
		endpoint += "/" + tf_app_id
		app_id = tf_app_id
	} else {
		return invokeErrorDiagsWithSummary(error_summary, "Either app_name or app_id should be given to add the new groups to an existing app.")
	}
	// new groups
	group_ids := d.Get("groups").([]interface{})
	if len(group_ids) > 0 {
		add_groups := make(map[string]interface{})
		for i := 0; i < len(group_ids); i++ {
			// default_permissions is in common.go
			add_groups[group_ids[i].(string)] = default_permissions
		}
		app_object := map[string]interface{}{
			"app_id": app_id,
			"add_groups": add_groups,
		}
		req, err := m.(*api_client).APICallBody("PATCH", endpoint, app_object)
		if err != nil {
			return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: PATCH %s: %v", endpoint, err))
		}
		d.SetId(req["app_id"].(string))
	} else {
		return invokeErrorDiagsWithSummary(error_summary, "Atleast one group should be given.")
	}
	return nil
}

// [R]: Read App
func resourceReadAppAssignGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Not required as it is not a actual resource
	return nil
}

// [U]: PATCH App to add new groups
func resourceUpdateAppAssignGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("groups") {
		old_group, new_group := d.GetChange("groups")
		// compute_add_and_del_arrays function is in common.go
		// Gets the new groups to be added(add_group_ids)
		add_group_ids, _ := compute_add_and_del_arrays(old_group, new_group)
		if len(add_group_ids) > 0 {
			endpoint := "sys/v1/apps/" + d.Id()
			add_groups := make(map[string]interface{})
			for i := 0; i < len(add_group_ids); i++ {
				// default_permissions is in common.go
				add_groups[add_group_ids[i]] = default_permissions
			}
			app_object := map[string]interface{}{
				"app_id": d.Id(),
				"add_groups": add_groups,
			}
			_, err := m.(*api_client).APICallBody("PATCH", endpoint, app_object)
			if err != nil {
				d.Set("groups", old_group)
				return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: PATCH %s: %v", endpoint, err))
			}
		}
	}

	return nil
}

// [D]: Delete App
func resourceDeleteAppAssignGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Simply deletes from the tf state
	d.SetId("")
	return nil
}

