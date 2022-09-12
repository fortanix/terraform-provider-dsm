# dsm\_app

## dsm\_app

Returns the Fortanix DSM App from the cluster as a resource 

## Usage Reference

```
resource "dsm_app" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description>
    new_credential = <true/false> 
```

## Update the groups

```

variable "app_data" {
    type = any
    default = {
        "mod_group" = {
            "<group_id>" : "<permissions(separated by '-')>",
            "<group_id>" : "<permissions(separated by '-')>"
        }
    }
}


resource "dsm_app" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description>
    patch_request  = true
    mod_group      = "${var.app_data.mod_groups}"
    del_group      = [<group_id>,<group_id>]
    other_group    = [<group_id>,<group_id>]
```

```

## Argument Reference

The following arguments are supported in the `dsm_app` resource block:

* **name**: The Fortanix DSM App name
* **default_group**: The Fortanix DSM group object id to be mapped to the app by default
* _**other_group (optional)**_: The Fortanix DSM group object id the app needs to be assigned to
* _**description (optional)**_: The description of the app 
* _**new\_credential (optional)**_: Set this if you want to rotate/regenerate the API key. The values can be set as `True`/`False`
* _**patch_request (optional)**_: Set this if you want to modify the default_group, to add the new groups, to delete the groups and to modify the permissions of groups. The values can be set as `True`/`False`
* _**del_groups (optional)**_: To delete the groups
* _**mod_groups (optional)**_: To modify the permissions of any group

   mod_groups example:

   A varaiable should be declared. Here it is named as app_data. Please follow the below varaible reference 
   to provide the permissions. For each group_id permissions should be given in a string format. Permissions
   are separated by '-'.

   mod_group     = "${var.app_data.mod_groups}"

   variable "app_data" {
    type = any
    default = {
        "mod_groups" = {
            "<group_id>" : "SIGN-VERIFY-ENCRYPT-DECRYPT-WRAPKEY-UNWRAPKEY-DERIVEKEY-MACGENERATE-MACVERIFY-EXPORT-MANAGE-AGREEKEY-AUDIT-TRANSFORM"
            "<group_id>" : "SIGN-ENCRYPT-DECRYPT-WRAPKEY-UNWRAPKEY-MACGENERATE-MACVERIFY-MANAGE-AGREEKEY-AUDIT-EXPORT"
        }
    }
   }


## Attribute Reference

The following attributes are stored in the `dsm_app` resource block:

* **id**: The unique ID of object from Terraform (matches the `app_id` from resource block)
* **name**: The App name from Fortanix DSM (matches the name provided during creation)
* **app\_id**: The unique ID of the app from Terraform
* **default\_group**: The default group name mapped to the Fortanix DSM app
* **acct\_id**: The account ID from Fortanix DSM
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The Fortanix DSM App description
* **credential**: The Fortanix DSM App API key
