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
    fpe             = <fpeOptions>
    description     = <sobject_description>
    custom_metadata = {        
                    <key> = <value>    
    }
    allowed_key_justifications_policy = <allowed_key_justifications_policy>
    allowed_missing_justifications = <allowed_missing_justifications>
    rsa             = <rsaOptions_string_format>
    elliptic_curve  = <elliptic_curve>
    value = <imported sobject content>
    hash_alg = <HashAlgorithm>
    subgroup_size = <subgroup_size_value>
    rotation_policy = {
      interval_days = <number of days>
      effective_at = "<yyyymmddThhmmssZ>"
      deactivate_rotated_key = <true/false>
      rotate_copied_keys = "all_external"
    }
}
```

## Argument Reference

The following arguments are supported in the `dsm_sobject` resource block:

* **name**: The security object name.
* **obj\_type**: The security object type.
* **key\_size**: The security object size. It should not be given only when the obj_type is EC.
* **group\_id**: The security object group assignment.
* _**key\_ops (optional)**_: The security object key permission.
* _**rsa (optional)**: The rsaOptions for an RSA object.
* _**description (optional)**_: The security object description.
* _**custom_metadata (optional)**_: The user defined security object attributes added to the key’s metadata from Fortanix DSM.
* _**fpe\_radix (optional)**_: integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).
* _**enabled (optional)**_: Whether the security object is enabled or disabled. The values are `True`/`False` 
* _**expiry date (optional)**_: The security object expiry date in RFC format.
* _**state (optional)**_: The state of the secret security object. Allowed states are: `None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`.
* _**rotate(optional)**_: specify method to use for key rotation.
  * **DSM** - To rotate from a DSM local key. The key material of new key will be stored in DSM.
* _**rotate_from(optional)**_ : Name of the security object to be rotated from
* _**elliptic_curve**_ : Standardized elliptic curve. It should be given only when the obj_type is EC.
* _**value**_  = Sobject content when importing content.
* _**allowed\_key\_justifications\_policy (optional)**_: The security object key justification policies for GCP External Key Manager.
* _**allowed\_missing\_justifications (optional)**_: The security object allows missing justifications even if not provided.
* _**hash\_alg**_ = Hashing Algorithm for KCDSA and ECKCDSA
* _**subgroup\_size**_ = Subgroup Size for DSA and ECKCDSA
* _**rotation_policy(optional)**_ = Policy to rotate a Security Object, configure the below parameters.
* * _**interval_days**_ = Rotate the key for every given number of days
* * _**interval_months**_ = Rotate the key for every given number of months
* * _**effective_at**_ = Start of the rotation policy time
* * _**deactivate_rotated_key**_ = Deactivate original key after rotation (true/false)
* * _**rotate_copied_keys**_ = Enable key rotation for copied keys
* _**fpe (optional)**_: FPE specific options. It should be given in string format like below:
```
This is a sample variable that specifies fpeOptions to create a Tokenization object that can tokenize credit card format data:

  variable "fpeOptionsExample" {
    type = any
    description = "The policy document. This is a JSON formatted string."
    default = <<-EOF
          { 
            "description": "Credit card",
            "format": {
              "char_set": [
                [
                  "0",
                  "9"
                ]
              ],
              "min_length": 13,
              "max_length": 19,
              "constraints": {
                "luhn_check": true
              }
          }
        }
        EOF
}

This is how we can reference this fpeOptions:
  fpe = var.fpeOptionsExample

Refer to the "fpeOptions" schema in https://www.fortanix.com/fortanix-restful-api-references/dsm for a better understanding of the fpe body.
```

## Note on rotational_policy

Only one of the following attributes should be used while configuring the interval in rotational_policy
  1. interval_days
  2. interval_months

## Attribute Reference

The following attributes are stored in the `dsm_sobject` resource block:

* **id**: The unique ID of object from Terraform (matches the `kid` from resource block).
* **kid**: The security object ID from Fortanix DSM.
* **name**: The security object name from Fortanix DSM (matches the name provided during creation).
*  **group_id**: The group object ID from Fortanix DSM.
* **acct\_id**: Account ID from Fortanix DSM.
* **obj\_type**: The security object key type from Fortanix DSM (matches the obj_type provided during creation).
* **key\_size**: The security object key size from Fortanix DSM (matches the key_size provided during creation).
* **key\_ops**: The security object key permission from Fortanix DSM.
  * Default is to allow all permissions except "EXPORT".
* **rsa**: rsaOptions passed as a string (if "RSA” `obj_type` is specified). The string should match the "rsa" value in Post body while working with Fortanix Rest API. For example, 
`rsa = "{\"encryption_policy\":[{\"padding\":{\"RAW_DECRYPT\":{}}},{\"padding\":{\"OAEP\":{\"mgf\":{\"mgf1\":{\"hash\":\"SHA1\"}}}}}],\"signature_policy\":[{\"padding\":{\"PKCS1_V15\":{}}},{\"padding\":{\"PSS\":{\"mgf\":{\"mgf1\":{\"hash\":\"SHA384\"}}}}}]}"`
* **creator**: The creator of the security object from Fortanix DSM.
  * **user**: If the security object was created by a user, the computed value will be the matching user id.
  * **app**: If the security object was created by a app, the computed value will be the matching app id.
* **description**: Security object description.
* **pub\_key**: Public key (if "RSA” `obj_type` is specified).
* **ssh\_pub\_key**: Open SSH public key (if "RSA” `obj_type` is specified).
* **state**: state of the secret (`None`, `PreActive`, `Active`, `Deactivated`, `Compromised`, `Destroyed`, `Deleted`).
* **expiry\_date**: The security object expiry date in RFC format.
* **custom\_metadata**: The user defined security object attributes added to the key’s metadata from Fortanix DSM.
* **fpe\_radix**:   integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).
* _**fpe (optional)**_: FPE specific options.
* **elliptic\_curve**: Standardized elliptic curve.
* _**allowed\_key\_justifications\_policy (optional)**_: The security object key justification policies for GCP External Key Manager. The allowed permissions are:  `CUSTOMER_INITIATED_SUPPORT` , `CUSTOMER_INITIATED_ACCESS`, `GOOGLE_INITIATED_SERVICE`, `GOOGLE_INITIATED_REVIEW`, `GOOGLE_INITIATED_SYSTEM_OPERATION`,  `THIRD_PARTY_DATA_REQUEST`,`REASON_NOT_EXPECTED`, `REASON_UNSPECIFIED`, `MODIFIED_CUSTOMER_INITIATED_ACCESS`, `MODIFIED_GOOGLE_INITIATED_SYSTEM_OPERATION`, `GOOGLE_RESPONSE_TO_PRODUCTION_ALERT`.
* _**allowed\_missing\_justifications (optional)**_: Boolean value which allows missing justifications even if not provided to the security object. The values are `True` / `False`.

* _**hash\_alg**_ = Hashing Algorithm for KCDSA and ECKCDSA. The allowed Hashing Algorithms are `SHA1`,`SHA224`, `SHA256`, `SHA384`, `SHA521`.
* _**subgroup\_size**_ = Subgroup Size for DSA and ECKCDSA. The allowed Subgroup Sizes are `224` and `256`
* _**rotation\_policy**_ = Policy to rotate a security object
