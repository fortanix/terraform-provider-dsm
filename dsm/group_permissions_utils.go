package dsm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var default_permissions = []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"}

// subjected to resource_app and resource_gcp_ekm_sa
// It will add and delete the groups if applicable
func getChangesInOtherGroups(d *schema.ResourceData, app_object map[string]interface{}) {
	old_group, new_group := d.GetChange("other_group")
	// compute_add_and_del_arrays function is in common.go
	add_group_ids, del_group_ids := compute_add_and_del_arrays(old_group, new_group)
	//Add the permissions to new groups
	add_groups := make(map[string]interface{})
	if len(add_group_ids) > 0 {
		add_group_perms := form_group_permissions(d.Get("other_group_permissions"))
		for i := 0; i < len(add_group_ids); i++ {
			if perms, ok := add_group_perms[add_group_ids[i]]; ok {
				add_groups[add_group_ids[i]] = perms
			} else {
				add_groups[add_group_ids[i]] = default_permissions
			}
		}
	}
	// Delete the groups
	if len(del_group_ids) > 0 {
		app_object["del_groups"] = del_group_ids
	}
	//Add the new groups
	if len(add_groups) > 0 {
		app_object["add_groups"] = add_groups
	}
}

// subjected to resource_app and resource_gcp_ekm_sa
// It adds the changes in group permissions
func getChangesInGroupPermissions(d *schema.ResourceData, app_object map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	mod_groups := make(map[string]interface{})
	for group_id, perms := range d.Get("mod_group_permissions").(map[string]interface{}) {
		mod_groups[group_id] = perms
	}
	if len(mod_groups) > 0 {
		app_mod_groups := make(map[string]interface{})
		//if default_group has changes in permissions
		default_group := d.Get("default_group").(string)
		if perms, ok := mod_groups[default_group]; ok {
			app_mod_groups[default_group] = strings.Split(perms.(string), ",")
			delete(mod_groups, default_group)
		}
		//checking whether all the group_ids from mod_group_permissions exists in other groups or not
		//if not it will ignore the mod_group_permissions of the unavailable group_id
		var other_group_latest []string
		if other_groups := d.Get("other_group").([]interface{}); len(other_groups) > 0 {
			for _, group_id := range other_groups {
				other_group_latest = append(other_group_latest, group_id.(string))
			}
		}
		for i := 0; i < len(other_group_latest); i++ {
			if perms, ok := mod_groups[other_group_latest[i]]; ok {
				app_mod_groups[other_group_latest[i]] = strings.Split(perms.(string), ",")
				delete(mod_groups, other_group_latest[i])
			}
		}
		if len(mod_groups) > 0 {
			var unavailable_group_ids []string
			for group_id := range mod_groups {
				unavailable_group_ids = append(unavailable_group_ids, group_id)
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "All the group_ids available in mod_group_permissions are not available in other_group.Please correct them.",
				Detail:   fmt.Sprintf("[E]: Input: mod_group_permissions: Please remove the group_ids from mod_group_permissions those are not part of other_group. \n Following group_ids are not available in other_group:\n %v", unavailable_group_ids),
			})
			if old_group, new_group := d.GetChange("other_group"); len(new_group.([]interface{})) > 0 {
				d.Set("other_group", old_group)
			}
			return diags
		}
		if len(app_mod_groups) > 0 {
			app_object["mod_groups"] = app_mod_groups
		}
	}
	return nil
}

// subjected to resource_app and resource_gcp_ekm_sa
// Forms new groups
func formAddGroups(d *schema.ResourceData, app_object map[string]interface{}) {
	add_group_perms := form_group_permissions(d.Get("other_group_permissions"))
	add_groups := make(map[string]interface{})
	if groups := d.Get("other_group").([]interface{}); len(groups) > 0 {
		for _, group_id := range groups {
			if perms, ok := add_group_perms[group_id.(string)]; ok {
				add_groups[group_id.(string)] = perms
			} else {
				add_groups[group_id.(string)] = default_permissions
			}
		}
	}
	if perms, ok := add_group_perms[d.Get("default_group").(string)]; ok {
		add_groups[d.Get("default_group").(string)] = perms
	} else {
		add_groups[d.Get("default_group").(string)] = default_permissions
	}
	app_object["add_groups"] = add_groups
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
