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

## Argument Reference

The following arguments are supported in the `dsm_app` resource block:

* **name**: The Fortanix DSM App name
* **default_group**: The Fortanix DSM group object id to be mapped to the app by default
* _**other_group (optional)**_: The Fortanix DSM group object id the app needs to be assigned to
* _**description (optional)**_: The description of the app 
* _**new\_credential (optional)**_: Set this if you want to rotate/regenerate the API key. The values can be set as `True`/`False`

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
