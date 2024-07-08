# dsm\_gcp\_sobject

## dsm\_gcp\_sobject

Returns the DSM security object from the cluster as a Resource for GCP CDC Group. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to GCP KMS as a Customer Managed Key (CMK).

## Usage Reference

```
resource "dsm_gcp_sobject" "sobject" {
    name            = <sobject_name>
    group_id        = <sobject_group_id>
    description     = <sobject_description>
    obj_type        = <key_type>
    key_size        = <key_size>
    key_ops         = <key_ops>
    enabled         = <true/false>
    expiry_date     = <expiry_date_RFC_format>
    key             = {
                    kid = <local_sobject_id> 
    } 
    custom_metadata = {
        gcp-key-id  = <gcp_key_name>
    }
    rotation_policy = {
      interval_days = <number_of_days>
      effective_at = "<yyyymmddThhmmssZ>"
      deactivate_rotated_key = <true/false>
    }
}
```

## Argument Reference

The following arguments are supported in the `dsm_gcp_sobject` resource block:

* **name**: The security object name
* **group\_id**: The GCP group ID in Fortanix DSM into which the key will be generated
* _**description (optional)**_: The security object description
* **key**: A Local security object imported to Fortanix DSM(BYOK) and copied to GCP KMS
* _**obj\_type (optional)**_: The type of security object
* _**key\_size (optional)**_: The size of the security object
* _**key\_ops (optional)**_: The security object operations permitted
* _**enabled (optional)**_: Whether the security object will be `Enabled` or `Disabled`. The values are `True`/`False`
* _**state (optional)**_: The key states of the GCP KMS key. The values are `Created`, `Deleted`, `Purged`
* _**expiry\_date (optional)**_: The security object expiry date in RFC format
* _**custom\_metadata (optional)**_:  GCP KMS Key metadata information
    *	**gcp-key-id** - Key name within GCP KMS
* _**rotation_policy(optional)**_ = Policy to rotate a Security Object, configure the below parameters.
* * _**interval_days**_ = Rotate the key for every given number of days
* * _**interval_months**_ = Rotate the key for every given number of months
* * _**effective_at**_ = Start of the rotation policy time
* * _**deactivate_rotated_key**_ = Deactivate original key after rotation (true/false)

## Note on rotational_policy

Only one of the following attributes should be used while configuring the interval in rotational_policy
1. interval_days
2. interval_months

## Attribute Reference

The following attributes are stored in the `dsm_gcp_sobject` resource block:

* **kid**: The security object ID from Fortanix DSM
* **name**: The security object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The account ID from Fortanix DSM
* **key**: A Local security object imported to Fortanix DSM(BYOK) and copied to GCP KMS
* **key\_ops**: The security object operations permitted from Fortanix DSM
    * Default is to copy all permissions from the local security object
* **links**: Link between local security object and GCP KMS security object
* **enabled**: Returns `True` or `False` if the security object is `Enabled` or `Disabled`
* **creator**: The creator of the security object from DSM
    * **user**: If the security object was created by a user, the computed value will be the matching user id
    * **app**: If the security object was created by a app, the computed value will be the matching app id
* **obj\_type**: The type of security object
* **key\_size**: The size of the security object
* **description**: The security object description
* **expiry\_date**: The security object expiry date in RFC format from Fortanix DSM
* _**custom\_metadata (optional)**_:  GCP KMS Key metadata information
    *	**gcp-key-id** – Key name within GCP KMS
* _**rotation\_policy**_ = Policy to rotate a Security Object
