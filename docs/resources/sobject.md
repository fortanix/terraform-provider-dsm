---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dsm_sobject Resource - terraform-provider-dsm"
subcategory: ""
description: |-
  Creates a new security object. The returned resource object contains the UUID of the security object for further references.
  A key value can be imported as a security object. This resource also can rotate or copy a security object.
---

# dsm_sobject (Resource)

Creates a new security object. The returned resource object contains the UUID of the security object for further references.
A key value can be imported as a security object. This resource also can rotate or copy a security object.

## Example Usage

```terraform
// Create a group
resource "dsm_group" "group" {
  name = "group"
  description = "group description"
}
// Create a security object in the above group
resource "dsm_sobject" "sobject" {
  name            = "sobject"
  obj_type        = "AES"
  group_id        = dsm_group.group.id
  key_size        = 256
  key_ops         = [
    "ENCRYPT",
    "DECRYPT",
    "WRAPKEY",
    "UNWRAPKEY",
    "DERIVEKEY",
    "MACGENERATE",
    "MACVERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "sobject description"
  custom_metadata = {
    key1 = "value1"
  }
  allowed_key_justifications_policy = [
    "CUSTOMER_INITIATED_SUPPORT",
    "CUSTOMER_INITIATED_ACCESS"
  ]
  allowed_missing_justifications = true
  rotation_policy = {
    interval_days = 20
    effective_at = "20231130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The security object group assignment.
- `name` (String) The security object name.
- `obj_type` (String) The security object type.
   * `Supported security objects`: AES, DES, DES3, RSA, DSA, KCDSA, EC, ECKCDSA, ARIA, SEED and Tokenization(fpe).

### Optional

- `allowed_key_justifications_policy` (List of String) The security object key justification policies for GCP External Key Manager. The allowed permissions are:
   * CUSTOMER_INITIATED_SUPPORT
   * CUSTOMER_INITIATED_ACCESS
   * GOOGLE_INITIATED_SERVICE
   * GOOGLE_INITIATED_REVIEW
   * GOOGLE_INITIATED_SYSTEM_OPERATION
   * THIRD_PARTY_DATA_REQUEST
   * REASON_NOT_EXPECTED
   * REASON_UNSPECIFIED
   * MODIFIED_CUSTOMER_INITIATED_ACCESS
   * MODIFIED_GOOGLE_INITIATED_SYSTEM_OPERATION
   * GOOGLE_RESPONSE_TO_PRODUCTION_ALERT
- `allowed_missing_justifications` (Boolean) Boolean value which allows missing justifications even if not provided to the security object. The values are True / False.
- `custom_metadata` (Map of String) The user defined security object attributes added to the key’s metadata from Fortanix DSM.
- `description` (String) The security object description.
- `elliptic_curve` (String) Standardized elliptic curve. It should be given only when the obj_type is EC or ECKCDSA.

| obj_type | Curve | key_ops |
| -------- | -------- |-------- |
| `EC` | SecP192K1, SecP224K1, SecP256K1  NistP192, NistP224, NistP256, NistP384, NistP521, X25519, Ed25519 | APPMANAGEABLE, SIGN, VERIFY, AGREEKEY, EXPORT |
| `ECKCDSA` | SecP192K1, SecP224K1, SecP256K1  NistP192, NistP224, NistP256, NistP384, NistP521 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |
- `enabled` (Boolean) Whether the security object is enabled or disabled.
   * The values are true/false.
- `expiry_date` (String) The security object expiry date in RFC format.
- `fpe` (String) FPE specific options. obj_type should be AES. It should be given in string format like below:
```This is a sample variable that specifies fpeOptions to create a Tokenization object that can tokenize credit card format data:
    variable "fpeOptionsExample" { 
      type = any
      description = "The policy document. This is a JSON formatted string."
      default = <<-EOF 
              {
               "description":"Credit card"
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

Refer to the fpeOptions schema in https://www.fortanix.com/fortanix-restful-api-references/dsm for a better understanding of the fpe body.
```
- `fpe_radix` (Number) integer, The base for input data. The radix should be a number from 2 to 36, inclusive. Each radix corresponds to a subset of ASCII alphanumeric characters (with all letters being uppercase). For instance, a radix of 10 corresponds to a character set consisting of the digits from 0 to 9, while a character set of 16 corresponds to a character set consisting of all hexadecimal digits (with letters A-F being uppercase).
- `hash_alg` (String) Hashing Algorithm for KCDSA and ECKCDSA.

| obj_type | hash_alg |
| -------- | -------- |
| `ECKCDSA` | SHA1,SHA224, SHA256, SHA384, SHA521|
| `KCDSA` | SHA224, SHA256 |
- `key_ops` (List of String) The security object key permission from Fortanix DSM.
   * Default is to allow all permissions except EXPORT
- `key_size` (Number) The security object size. It should not be given only when the obj_type is EC and ECKCDSA.

| obj_type | key_size | key_ops |
| -------- | -------- |-------- |
| `RSA` | 1024, 2048, 4096, 8192 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |
| `DSA` | 2048, 3072 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |
| `KCDSA` | 2048 | APPMANAGEABLE, SIGN, VERIFY, EXPORT |
| `AES` | 128, 192, 256 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |
| `DES` | 56 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, APPMANAGEABLE, EXPORT |
| `DES3` | 112, 168 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |
| `ARIA` | 128, 192, 256 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |
| `SEED` | 128 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, EXPORT |
- `rotate` (String) specify method to use for key rotation.
- `rotate_from` (String) Name of the security object to be rotated from.
- `rotation_policy` (Map of String) Policy to rotate a Security Object, configure the below parameters.
   * `interval_days`: Rotate the key for every given number of days.
   * `interval_months`: Rotate the key for every given number of months.
   * `effective_at`: Start of the rotation policy time.
   * `rotate_copied_keys`: Enable key rotation for copied keys.
   * `deactivate_rotated_key`: Deactivate original key after rotation true/false.
   * **Note:** Either interval_days or interval_months should be given, but not both.
- `rsa` (String) rsaOptions passed as a string (if ”RSA” obj_type is specified). The string should match the 'rsa' value in Post body while working with Fortanix Rest API. For Example:

`rsa = "{\"encryption_policy\":[{\"padding\":{\"RAW_DECRYPT\":{}}},{\"padding\":{\"OAEP\":{\"mgf\":{\"mgf1\":{\"hash\":\"SHA1\"}}}}}],\"signature_policy\":[{\"padding\":{\"PKCS1_V15\":{}}},{\"padding\":{\"PSS\":{\"mgf\":{\"mgf1\":{\"hash\":\"SHA384\"}}}}}]}"`
- `state` (String) The state of the secret security object.
   * Allowed states are: None, PreActive, Active, Deactivated, Compromised, Destroyed, Deleted.
- `subgroup_size` (Number) Subgroup Size for DSA and ECKCDSA. The allowed Subgroup Sizes are 224 and 256.

| obj_type | subgroup_size | usage
| -------- | -------- | -------- |
| `DSA` | 224, 256| 224: When DSA key_size is 2048. 256: When DSA key_size is 2048 and 3072.
| `KCDSA` | 224, 256| 224, 256: When KCDSA key_size is 2048.
- `value` (String) Sobject content when importing content.

### Read-Only

- `acct_id` (String) Account ID from Fortanix DSM.
- `copied_from` (String) Security object that is copied to the current security object.
- `copied_to` (List of String) List of security objects copied by the current security object.
- `creator` (Map of String) The creator of the security object from Fortanix DSM.
   * `user`: If the security object was created by a user, the computed value will be the matching user id.
   * `app`: If the security object was created by a app, the computed value will be the matching app id.
- `dsm_name` (String) The security object name.
- `id` (String) The ID of this resource.
- `kid` (String) The security object ID from Fortanix DSM.
- `pub_key` (String) Public key (if ”RSA” obj_type is specified).
- `replaced` (String) Replaced by a security object.
- `replacement` (String) Replacement of a security object.
- `ssh_pub_key` (String) Open SSH public key (if ”RSA” obj_type is specified).
