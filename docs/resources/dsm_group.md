# dsm\_group

## dsm\_group

Returns the Fortanix DSM group object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_group" "group" {
    name = <group_name>
    description = <group_description>
    approval_policy = <group_quorum_policy_definition>
    hmg = <group_HMS/KMS_definition>
    key_undo_policy_window_time =<key_undo_policy_window_time>
}
```

## Argument Reference

The following arguments are supported in the `dsm_group` resource block:

* **name**: The Fortanix DSM group object name.
* _**description (optional)**_: The Fortanix DSM group object description
* _**approval_policy (optional)**_: The Fortanix DSM group object quorum approval policy definition as a JSON string
* _**hmg (optional)**_: The Fortanix DSM group object HMS/KMS definition as a JSON string
* _**key_undo_policy_window_time(optional)**_: The Fortanix DSM group object key undo policy window time as an Integer(Number of days).

## Attribute Reference

The following attributes are stored in the `dsm_group` resource block:

* **id**: Unique ID of object from Terraform (matches the `group_id` from resource block)
* **group\_id**: Group object ID from Fortanix DSM
* **name**: Group object name from Fortanix DSM (matches the `name` provided during creation)
* **acct\_id**: Account ID from Fortanix DSM
* **creator**: Creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The Fortanix DSM group object description
* **approval_policy**: The Fortanix DSM group object quorum approval policy definition as a JSON string
* **hmg**: The Fortanix DSM group object HMS/KMS definition as a JSON string
* _**key_undo_policy_window_time**_: The Fortanix DSM group object key undo policy window time as an Integer(Number of days).
