# dsm\_group\_crypto\_policy

## dsm\_group\_crypto\_policy

Returns the Fortanix DSM group cryptographic policy object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_group_crypto_policy" "name" {
    name = <group_name>
    cryptographic_policy = <group_cryptographic_policy_definition>
    <depends_on = [resource.dsm_group.GROUPNAME]> // optional but recommended
}
```

## Argument Reference

The following arguments are supported and required in the `dsm_group_crypto_policy` resource block:

* **name**: The Fortanix DSM group object name.
* **cryptographic\_policy**: The Fortanix DSM group object cryptographic policy definition as a JSON string

## Attribute Reference

The following attributes are stored in the `dsm_group_crypto_policy` resource block:

* **id**: Unique ID of object from Terraform (matches the `group_id` from resource block)
* **group\_id**: Group object ID from Fortanix DSM
* **name**: Group object name from Fortanix DSM (matches the `name` provided during creation)
* **acct\_id**: Account ID from Fortanix DSM
* **creator**: Creator of the group object from Fortanix DSM
* **description**: The Fortanix DSM group object description
* **approval\_policy**: The Fortanix DSM group object quorum approval policy definition as a JSON string
* **cryptographic\_policy**: The Fortanix DSM group object cryptographic policy definition as a JSON string

## Important Note
It is best to use a "depends_on" directive to wait for the creation of the parent group resource:
depends_on = [resource.dsm_group.GROUPNAME]