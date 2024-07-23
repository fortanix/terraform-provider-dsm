# dsm\_group

## dsm\_group

Returns the Fortanix DSM group object from the cluster as a Data Source.

## Usage Reference

```
data "dsm_group" "group" {
    name = <group_name>
}
```

## Argument Reference

The following arguments are supported in the `dsm_group` data block:

* **name**: Group object name

## Attribute Reference

The following attributes are stored in the `dsm_group` data source block:

* **id**: The unique ID of object from Terraform (matches the `group_id`) 
* **group\_id**: The group object ID from Fortanix DSM
* **name**: The group object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The account ID from Fortanix DSM
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The group object description from Fortanix DSM
* **approval_policy**: The Fortanix DSM group object quorum approval policy definition as a JSON string
* **cryptographic_policy**: The Fortanix DSM group object cryptographic policy definition as a JSON string
* **hmg**: The Fortanix DSM group object HMS/KMS definition as a JSON string