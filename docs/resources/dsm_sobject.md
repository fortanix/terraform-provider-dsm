# dsm\_sobject

## dsm\_sobject

Returns the Fortanix DSM security object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_sobject" "sobject" {
    name            = <sobject_name>
    obj_type        = <key_type>
    group_id        = <sobject_group_id>
    key_size        = <sobject_key_size>
    key_ops         = <key_ops>
    enabled         = <true/false>
    expiry_date     = <expiry_date_RFC_format>
    fpe_radix       = <fpe_radix>
    description     = <sobject_description>
    key             = {        
    kid             = <local_sobject_id>     
    }     
    custom_metadata = {        
                    <key> = <value>    
    }
}
```

## Argument Reference

The following arguments are supported in the `dsm_sobject` resource block:

* **name**: The security object name
* **obj\_type**: The security object type
* **key\_size**: The security object size
* **group\_id**: The security object group assignment
* _**key\_ops (optional)**_: The security object key permission
* _**description (optional)**_: The security object description
* _**custom_metadata (optional)**_: The user defined security object attributes added to the key’s metadata from Fortanix DSM
* _**fpe\_radix (optional)**_: integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).
* _**enabled (optional)**_: Whether the security object is enabled or disabled. The values are `True`/`False` 
* _**expiry date (optional)**_: The security object expiry date in RFC format 
* _**state (optional)**_: The state of the secret security object. Allowed states are: `None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`
* _**rotate(optional)**_: specify method to use for key rotation 
  * **DSM** - To rotate from a DSM local key. The key material of new key will be stored in DSM.
* _**rotate_from(optional)**_  = Name of the security object to be rotated from

## Attribute Reference

The following attributes are stored in the `dsm_sobject` resource block:

* **id**: The unique ID of object from Terraform (matches the `kid` from resource block)
* **kid**: The security object ID from Fortanix DSM
* **name**: The security object name from Fortanix DSM (matches the name provided during creation)
*  **group_id**: The group object ID from Fortanix DSM
* **acct\_id**: Account ID from Fortanix DSM
* **obj\_type**: The security object key type from Fortanix DSM (matches the obj_type provided during creation)
* **key\_size**: The security object key size from Fortanix DSM (matches the key_size provided during creation)
* **key\_ops**: The security object key permission from Fortanix DSM
  * Default is to allow all permissions except "EXPORT"
* **creator**: The creator of the security object from Fortanix DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **description**: Security object description
* **ssh\_pub\_key**: Open SSH public key (if "RSA” `obj_type` is specified)
* **state**: state of the secret (`None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`)
* **expiry\_date**: The security object expiry date in RFC format
* **custom\_metadata**: The user defined security object attributes added to the key’s metadata from Fortanix DSM.
* **fpe\_radix**:   integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).
