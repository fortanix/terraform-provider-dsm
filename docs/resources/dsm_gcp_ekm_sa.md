# dsm\_gcp\_ekm\_sa

## dsm\_gcp\_ekm\_sa

Returns the Fortanix DSM Google EKM app from the cluster as a Resource.

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


resource "dsm_app" "gcp_ekm_sa" {
    name                    = <google_service_account_name>
    default_group           = <DSM group id>
    description             = <description of the app>
    other_group             = [<group_id>,<group_id>]
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


resource "dsm_app" "gcp_ekm_sa" {
    name                    = <google_service_account_name>
    default_group           = <DSM group id>
    description             = <description of the app>
    other_group             = [<group_id>,<group_id>]
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
    
    /* add the new groups:  just add the new group_ids in other_group array.
     * delete the existing groups: just remove group_id from other_group array.
     */
}
```


## Argument Reference

The following arguments are supported in the `dsm_gcp_ekm_sa` resource block:

* **name**: The Google service account name
* **default\_group**: The Fortanix DSM group id to be mapped to the app by default
* _**description (optional)**_: The description of the app
* _**other_group (optional)**_: The Fortanix DSM group object id the app needs to be assigned to. If you want to
                                delete the existing groups from an app, remove the ids during update.
* _**other_group_permissions(optional)**_: Incase if you want to change the default permissions of a new group.
* _**mod_group_permissions (optional)**_: To modify the permissions of any existing group

```
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
```


## Attribute Reference

The following attributes are stored in the `dsm_gcp_ekm_sa` resource block:

* **id**: The unique ID of object from Terraform (matches the `app_id` from resource block)
* **name**: The Google service account name
* **app\_id**: The unique ID of app from Terraform
* **default\_group**: The default group id mapped to the app
* **acct\_id**: The account ID from Fortanix DSM
* **creator**: The creator of the security object from Fortanix DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **description**: The Fortanix DSM App description
