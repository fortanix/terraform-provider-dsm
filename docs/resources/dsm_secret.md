# dsm\_secret

## dsm\_secret

Returns the Fortanix DSM secret security object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_secret" "secret" {
    name            = <secret_name>
    group_id        = <secret_group_id>
    description     = <sobject_description>
    enabled         = <true/false>
    state           = <secret_state>
    value           = <secret_value>
    expiry_date     = <expiry_date_RFC_format>
    custom_metadata = {        
                    <key>  = <value>    
    }
}
```

## Argument Reference

The following arguments are supported in the `dsm_secret` resource block:

* **name**: The Fortanix DSM secret security object name
* **group\_id**: The Fortanix DSM security object group assignment
* _**custom\_metadata (optional)**_: The user defined security object attributes added to the key’s metadata
* _**enabled (optional)**_: Whether the security object is `Enabled` or `Disabled`. The values are `True`/`False`
* _**expiry date (optional)**_: The security object expiry date in RFC format 
* _**value (optional)**_: The secret value
* _**state (optional)**_: The state of the secret security object. Allowed states are: `None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`
* _**description (optional)**_: The Fortanix DSM security object description
* _**rotate(optional)**_: boolean value true/false to enable/disable rotation 
* _**rotate_from(optional)**_  = Name of the security object to be rotated from

## Attribute Reference

The following attributes are stored in the `dsm_secret` resource block:

* **id**: The unique ID of object from Terraform (matches the `kid` from resource block)
* **kid**: Security object ID from Fortanix DSM
* **name**: Security object name from Fortanix DSM (matches the `name` provided during creation)
* **acct\_id**: Account ID from Fortanix DSM
* **group\_id**: The group object ID from Fortanix DSM
* **obj\_type**: The security object key type from Fortanix DSM (matches the obj_type provided during creation)
* **key\_ops**: The security object key permission from Fortanix DSM
  * Default is to allow all permissions except "EXPORT"
* **enabled**: Returns `True` or `False` if the security object is enabled or disabled
* **creator**: Creator of the security object from Fortanix DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **description**: Security object description
* **state**: The state of the secret security object from Fortanix DSM (`None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`)
* **expiry\_date**: The secret security object expiry date in RFC format
* **custom\_metadata**: The user defined security object attributes added to the key’s metadata from Fortanix DSM

