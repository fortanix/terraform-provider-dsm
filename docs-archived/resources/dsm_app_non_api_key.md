# dsm\_app_non_api_key

## dsm\_app_non_api_key

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
```
How to create an AWS XKS app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "awsxks"
    }
}
```
How to create an AWS IAM app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "awsiam"
    }
}
```
How to create a Certificate app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "certificate"
        certificate = "<certificate_value>"
    }
}
```
How to create a Trusted CA app?
```
Example of IP address

resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "trustedca"
        ca_certificate = "<certificate_value>"
        ip_address = "<ip_address>"
    }
}

Example of DNS name

resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "trustedca"
        ca_certificate = "<certificate_value>"
        dns_name = "<dns_name>"
    }
}
```


## Update the App

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


resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    description    = <app_description>
    other_group    = [<group_id>,<group_id>] 
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions 
}

```

How to create an AWS XKS app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    authentication_method = {
        type = "awsxks"
    }
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
}
```
How to create an AWS IAM app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "awsiam"
    }
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
}
```
How to create a Certificate app?
```
resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "certificate"
        certificate = "<certificate_value>"
    }
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
}
```
How to create a Trusted CA app?
```
Example of IP address

resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "trustedca"
        ca_certificate = "<certificate_value>"
        ip_address = "<ip_address>"
    }
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
}

Example of DNS name

resource "dsm_app_non_api_key" "app" {
    name           = <app_name>
    default_group  = <group_id>
    other_group    = [<group_id>,<group_id>]
    description    = <app_description> 
    other_group_permissions = local.app_other_group_permissions
    authentication_method = {
        type = "trustedca"
        ca_certificate = "<certificate_value>"
        dns_name = "<dns_name>"
    }
    /* add the new groups:  just add the new group_ids in this array.
     * delete the existing groups: just remove group_id from this array.
     */
    other_group_permissions = local.app_other_group_permissions
    mod_group_permissions   = local.app_mod_group_permissions
}
```


## Argument Reference

The following arguments are supported in the `dsm_app_non_api_key` resource block:

* **name**: The Fortanix DSM App name
* **default_group**: The Fortanix DSM group object id to be mapped to the app by default
* _**other_group (optional)**_: The Fortanix DSM group object id the app needs to be assigned to. If you want to
  delete the existing groups from an app, remove the ids during update.
* _**description (optional)**_: The description of the app
* _**other_group_permissions(optional)**_: Incase if you want to change the default permissions of a new group.
* _**mod_group_permissions (optional)**_: To modify the permissions of any existing group
* * _**authentication_method**_: To modify the permissions of any existing group
  * _**type**_: Following authentication types are supported.
    * awsxks
    * awsiam
    * certificate
    * trustedca
  * _**certificate**_: Certificate value, this should be configured when the type is certificate
  * _**ca_certificate**_: CA certificate value, this should be configure when the type is trustedca
    * One of the following values should be configured when the type is trustedca
    * _**ip_address**_: IP address value for trusted ca
    * _**dns_name**_: DNS name for trusted ca

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

The following attributes are stored in the `dsm_app_non_api_key` resource block:

* **id**: The unique ID of object from Terraform (matches the `app_id` from resource block)
* **name**: The App name from Fortanix DSM (matches the name provided during creation)
* **app\_id**: The unique ID of the app from Terraform
* **default\_group**: The default group name mapped to the Fortanix DSM app
* **acct\_id**: The account ID from Fortanix DSM
* **creator**: The creator of the group object from Fortanix DSM
    * **user**: If the group was created by a user, the computed value will be the matching user id
    * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The Fortanix DSM App description
* **credential**: The Fortanix DSM App credential, AWS xks access and secret key will be stored
* **authentication_method**_: The  Authentication method details