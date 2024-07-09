# dsm\_group\_user\_role

## dsm\_group\_user\_role

Associates a user with a group and a role.

## Usage Reference

```
resource "dsm_group_user_role" "binding_name" {
    group_name = <group_name>
    user_email = <user_email>
    role_name = <role_name>
}
```

## Argument Reference

The following arguments are supported and required in the `dsm_group_user_role` resource block:

* **group\_name**: The Fortanix DSM group object name.
* **user\_email**: The Fortanix DSM user object email.
* **role\_name**: The Fortanix DSM role object name.

## Attribute Reference

The following attributes are stored in the `dsm_group_user_role` resource block:

* **id**: Unique ID of object from Terraform (matches the `user_id` from resource block)
* **group\_name**: Group object name from Fortanix DSM (matches the `group_name` provided during creation)
* **user\_email**: User object email from Fortanix DSM (matches the `user_email` provided during creation)
* **role\_name**: Role object name from Fortanix DSM (matches the `role_name` provided during creation)
* **user\_id**: User object ID from Fortanix DSM
* **group\_id**: Group object ID from Fortanix DSM
* **role\_id**: Role object ID from Fortanix DSM

## Important Note
If applicable, it is best to use a "depends_on" directive to wait for the creation of a parent resource:</br>
depends_on = [resource.dsm_group.GROUPNAME]