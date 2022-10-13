# dsm\_app

## dsm\_app

Returns the Fortanix DSM App from the cluster as a resource 

## Usage Reference



```

locals {
    app_other_group_permissions = zipmap(
        [
            dsm_group.group1.group_id,
            "<group_id>"
        ],
        [
            "SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM",
            "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM"
        ]
    )
}


resource "dsm_app" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description>
    new_credential = <true/false> 
    other_group_permissions = local.app_other_group_permissions
}
```

## Update the groups

```

locals {
    app_mod_group_permissions = zipmap(
        [
            dsm_group.group1.group_id,
            "<group_id>"
        ],
        [
            "SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM",
            "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM"
        ]
    )
}


resource "dsm_app" "app" {
    name           = <app_name>
    default_group  = <group_id>
    description    = <app_description>
    other_group    = [<group_id>,<group_id>] 
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    new_credential = <true/false> 
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions 
}
```

```

## Argument Reference

The following arguments are supported in the `dsm_app` resource block:

* **name**: The Fortanix DSM App name
* **default_group**: The Fortanix DSM group object id to be mapped to the app by default
* _**other_group (optional)**_: The Fortanix DSM group object id the app needs to be assigned to. If you want to 
                                delete the existing groups from an app, remove the ids during update.
* _**description (optional)**_: The description of the app 
* _**new\_credential (optional)**_: Set this if you want to rotate/regenerate the API key. The values can be set as `True`/`False`
* _**other_group_permissions(optional)**_: Incase if you want to change the default permissions of a new group.
* _**mod_group_permissions (optional)**_: To modify the permissions of any existing group


   other_group_permissions example:
   
   A variable should be declared as locals. Here it is named as app_other_group_permissions. Please follow the below 
   varaible reference to provide the permissions. For each group_id permissions should be given in a string format. Permissions
   are separated by comma(","). Count of group_ids and permission strings should be same.
   First group_id in the first array will match to first string in the second array and so on.

   other_group_permissions = local.app_other_group_permissions
   
   locals {
        app_other_group_permissions = zipmap(
            [
                dsm_group.group1.group_id,
                "<group_id>"
            ],
            [
                "SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM",
                "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM"
            ]
        )
   }

   mod_group_permissions example:

   A variable should be declared as locals. Here it is named as app_other_group_permissions. Please follow the below 
   varaible reference to provide the permissions. For each group_id permissions should be given in a string format. 
   Permissions are separated by comma(","). Count of group_ids and permission strings should be same.
   First group_id in the first array will match to first string in the second array and so on.

   locals {
        app_mod_group_permissions = zipmap(
            [
                dsm_group.group1.group_id,
                "<group_id>"
            ],
            [
                "SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM",
                "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM"
            ]
        )
   }
   
   mod_group_permissions = local.app_mod_group_permissions



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
