# dsm\_secret

## dsm\_secret

Returns the Fortanix DSM secret object from the cluster as a Data Source.

## Usage Reference

```
data "dsm_secret" "secret" {
    name   = <secret_name>
    export = <true/false>
}
```

## Argument Reference

The following arguments are supported in the `dsm_group` resource block:

* **name**: The secret security object name in Fortanix DSM
* _**export (optional)**_: Exports the secret based on the value shown. The value is either `True`/`False`

## Attribute Reference

The following attributes are stored in the `dsm_group` data source block:

* **id**: The unique ID of object from Terraform (matches the `kid`) 
* **kid**: The unique ID of the secret from Fortanix DSM
* **group\_id**: The group object ID from Fortanix DSM
* **name**: The secret object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The account ID from Fortanix DSM
* **pub\_key**: Public key from DSM (If applicable)
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The group object description
* **value**: The (sensitive) value of the secret shown if exported in base64 format
