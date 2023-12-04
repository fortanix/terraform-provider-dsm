# dsm\_aws\_sobject

## dsm\_aws\_sobject

Returns the Fortanix DSM security object from the cluster as a Resource for AWS KMS Group. This is a Bring-Your-Own-Key (BYOK) method and copies an existing Fortanix DSM local security object to AWS KMS as a Customer Managed Key (CMK).

## Usage Reference

```
resource "dsm_aws_sobject" "sobject" {
    name                   = <sobject_name>
    group_id               = <sobject_group_id>
    description            = <sobject_description>
    obj_type               = <key_type>
    key_size               = <key_size>
    key_ops                = <key_ops>
    enabled                = <true/false>
    state                  = <aws_key_state>
    pending_window_in_days = <pending_deletion_window_default_7>
    expiry_date            = <expiry_date_RFC_format>
    key                    = {
                           kid = <local_sobject_id> 
    } 
    custom_metadata        = {
        aws-aliases        = <alias-to-use>
    }
    rotation_policy = {
      interval_days = <number of days>
      effective_at = "<yyyymmddThhmmssZ>"
    }
}
```

## Argument Reference

The following arguments are supported in the `dsm_aws_sobject` resource block:

* **name**: The security object name
* **group\_id**: The security object group assignment
* _**description (optional)**_: The security object description
* **key**: A Local security object imported to Fortanix DSM(BYOK) and copied to AWS KMS
* _**obj\_type (optional)**_: The type of security object
* _**key\_size (optional)**_: The size of the security object
* _**key\_ops (optional)**_: The security object operations permitted
* _**enabled (optional)**_: Whether the security object will be enabled or disabled. The values are `True`/`False`
* _**state (optional)**_: The key states of the AWS key. The values are `PendingDeletion`, `Enabled`, `Disabled`, `PendingImport`
* _**pending_window\_in\_days (optional)**_: The default value is `7` days, input the value for “`days`” after which the AWS key will be deleted 
* _**expiry\_date (optional)**_: The security object expiry date in RFC format
* _**rotate(optional)**_: specify method to use for key rotation 
  * **DSM** - To rotate from a DSM local key. The key material of new key will be stored in DSM.
  * **AWS** - To rotate from a AWS key. The key material of new key will be stored in AWS.
* _**rotate_from(optional)**_  = Name of the security object to be rotated

* _**custom\_metadata (optional)**_:  Contains metadata about an AWS KMS key
  *	**aws-aliases** – The display name for AWS KMS key used to identify the key.
  *	**aws-policy** - JSON format of AWS policy that should be enforced for the key.
* * _**rotation_policy(optional)**_ = Policy to rotate a Security Object, configure the below parameters.
* * _**interval_days**_ = Rotate the key for every given number of days
* * _**interval_weeks**_ = Rotate the key for every given number of weeks
* * _**interval_months**_ = Rotate the key for every given number of months
* * _**interval_years**_ = Rotate the key for every given number of years
* * _**effective_at**_ = Start of the rotation policy time

## Note on rotational_policy

Only one of the following attributes should be used while configuring the interval in rotational_policy
1. interval_days
2. interval_weeks
3. interval_months
4. interval_years


## Attribute Reference

The following attributes are stored in the `dsm_aws_sobject` resource block:

* **kid**: The security object ID from Fortanix DSM
* **name**: The security object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The account ID from Fortanix DSM
* **key**: A Local security object imported to Fortanix DSM(BYOK) and copied to AWS KMS
* **key\_ops**: The security object operations permitted from Fortanix DSM
  * Default is to copy all permissions from the local security object
* **links**: Link between local security object and AWS KMS security object
* **enabled**: true or false
* **creator**: The creator of the security object from DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **external**: AWS CMK level metadata 
  *	Key\_arn
  * Key\_id
  * Key\_state
  * Key\_aliases
  * Key\_deletion_date
* **obj\_type**: The type of security object 
* **key\_size**: The size of the security object
* **description**: The security object description
* **expiry\_date**: The security object expiry date in RFC format from Fortanix DSM
* * _**rotation_\_policy**_ = Policy to rotate a Security Object
