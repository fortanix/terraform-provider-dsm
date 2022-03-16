# dsm\_gcp\_ekm\_sa

## dsm\_gcp\_ekm\_sa

Returns the Fortanix DSM Google EKM app from the cluster as a Resource.

## Usage Reference

```
resource "dsm_ gcp_ekm_sa" "gcp_ekm_sa" {
    name          = <google_service_account_name> 
    default_group = <DSM group Name> 
    description   = <description of the app>
}
```

## Argument Reference

The following arguments are supported in the `dsm_gcp_ekm_sa` resource block:

* **name**: The Google service account name
* **default\_group**: The Fortanix DSM group name to be mapped to the app by default 
* _**description (optional)**_: The description of the app 

## Attribute Reference

The following attributes are stored in the `dsm_gcp_ekm_sa` resource block:

* **id**: The unique ID of object from Terraform (matches the `app_id` from resource block)
* **name**: The Google service account name
* **app\_id**: The unique ID of app from Terraform
* **default\_group**: The default group name mapped to the app
* **acct\_id**: The account ID from Fortanix DSM
* **creator**: The creator of the security object from Fortanix DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **description**: The Fortanix DSM App description
