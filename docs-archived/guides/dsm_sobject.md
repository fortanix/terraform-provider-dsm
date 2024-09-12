## Create an AES security object

```terraform

resource "dsm_group" "group" {
  name = "group"
}
resource "dsm_sobject" "aes_sobject_example" {
  name            = "aes_sobject_example"
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
  description     = "aes sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a DES security object

```terraform
resource "dsm_sobject" "des_sobject_example" {
  name            = "des_sobject_example"
  obj_type        = "DES"
  group_id        = dsm_group.group.id
  key_size        = 56
  key_ops         = [
    "ENCRYPT",
    "DECRYPT",
    "WRAPKEY",
    "UNWRAPKEY",
    "DERIVEKEY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "des sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a DES3 security object

```terraform
resource "dsm_sobject" "des3_sobject_example" {
  name            = "des3_sobject_example"
  obj_type        = "DES3"
  group_id        = dsm_group.group.id
  key_size        = 112
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
  description     = "des3 sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a RSA security object

```terraform
resource "dsm_sobject" "rsa_sobject_example" {
  name            = "rsa_sobject_example"
  obj_type        = "RSA"
  group_id        = dsm_group.group.id
  key_size        = 2048
  key_ops         = [
    "ENCRYPT",
    "DECRYPT",
    "WRAPKEY",
    "UNWRAPKEY",
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "rsa sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a DSA security object

```terraform
resource "dsm_sobject" "dsa_sobject_example" {
  name            = "dsa_sobject_example"
  obj_type        = "DSA"
  group_id        = dsm_group.group.id
  key_size        = 2048
  subgroup_size = 224
  key_ops         = [
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "dsa sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a EC security object

```terraform
resource "dsm_sobject" "ec_sobject_example" {
  name            = "ec_sobject_example"
  obj_type        = "EC"
  group_id        = dsm_group.group.id
  elliptic_curve  = "SecP192K1"
  key_ops         = [
    "SIGN",
    "VERIFY",
    "AGREEKEY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "EC sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a KCDSA security object

```terraform
resource "dsm_sobject" "kcdsa_sobject_example" {
  name            = "kcdsa_sobject_example"
  obj_type        = "KCDSA"
  group_id        = dsm_group.group.id
  key_size        = 2048
  subgroup_size   = 256
  hash_alg = "SHA256"
  key_ops         = [
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "KCDSA sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a EC-KCDSA security object

```terraform
resource "dsm_sobject" "eckcdsa_sobject_example" {
  name            = "eckcdsa_sobject_example"
  obj_type        = "ECKCDSA"
  group_id        = dsm_group.group.id
  elliptic_curve  = "SecP192K1"
  hash_alg = "SHA1"
  key_ops         = [
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "ECKCDSA sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create an ARIA security object

```terraform
resource "dsm_sobject" "aria_sobject_example" {
  name            = "aria_sobject_example"
  obj_type        = "ARIA"
  group_id        = dsm_group.group.id
  key_size        = 128
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
  description     = "ARIA sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a SEED security object

```terraform
resource "dsm_sobject" "seed_sobject_example" {
  name            = "seed_sobject_example"
  obj_type        = "SEED"
  group_id        = dsm_group.group.id
  key_size        = 128
  key_ops         = [
    "ENCRYPT",
    "DECRYPT",
    "WRAPKEY",
    "UNWRAPKEY",
    "DERIVEKEY",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "SEED sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a BLS security object

```terraform
resource "dsm_sobject" "bls_sobject_example" {
  name            = "bls_sobject_example"
  obj_type        = "BLS"
  group_id        = dsm_group.group.id
  bls = {
    variant = "small_signatures"
  }
  key_ops         = [
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  description     = "BLS sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a LMS security object

```terraform
resource "dsm_sobject" "lms_sobject_example" {
  name            = "lms_sobject_example"
  obj_type        = "LMS"
  group_id        = dsm_group.group.id
  lms = {
    l1_height = 5
    l2_height = 10
    node_size = 32
  }
  key_ops         = [
    "SIGN",
    "VERIFY",
    "APPMANAGEABLE"
  ]
  enabled         = true
  //expiry_date     = "2025-02-02T17:04:05Z"
  description     = "LMS sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}
```

## Create a Tokenization security object

```terraform
resource "dsm_sobject" "tokenization_sobject_example" {
  name            = "tokenization_sobject_example"
  group_id        = dsm_group.group1_example.id
  obj_type        = "AES"
  key_size        = 256
  fpe             = var.fpeOptionsExample
  key_ops         = [
    "ENCRYPT",
    "DECRYPT",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}

variable "fpeOptionsExample" {
  type = any
  description = "The policy document. This is a JSON formatted string."
  default = <<-EOF
              {
               "description":"Credit card",
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
```

## Import a security object (certificate)

```terraform
/*
To import any type of security object, value should be provided.
*/
resource "dsm_sobject" "certificate" {
  name            = "certificate_creation"
  obj_type        = "CERTIFICATE"
  group_id        = dsm_group.group.id
  value           = "XXXXXXXXXXXX<CERTIFICATE_value_in_a_string_format>XXXXXXXXXXXXXX"
  expiry_date     = "2025-02-02T17:04:05Z"
  enabled         = true
  key_ops         = [
    "ENCRYPT",
    "VERIFY",
    "WRAPKEY",
    "APPMANAGEABLE",
    "EXPORT"
  ]
  allowed_key_justifications_policy = [
    "CUSTOMER_INITIATED_SUPPORT",
    "CUSTOMER_INITIATED_ACCESS"
  ]
  allowed_missing_justifications = true
  custom_metadata = {
    key1 = "value1"
  }
}
```

## Copy a Security Object

```terraform
resource "dsm_sobject" "aes_sobject_example" {
  name            = "aes_sobject_example"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}

// Copy a security object in a group
resource "dsm_sobject" "aes_sobject_example_copy" {
  name            = "aes_sobject_example_copy"
  group_id        = dsm_group.group.id
  key             = {
    kid = dsm_sobject.aes_sobject_example.id
  }
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
  expiry_date     = "2026-02-02T17:04:05Z"
  description     = "sobject description copy"
  custom_metadata = {
    key1 = "value1"
  }
  rotation_policy = {
    interval_days = 20
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
}

```

## Rotate a security object

```terraform
resource "dsm_sobject" "aes_sobject_example_rotate" {
  name        = "aes_sobject_example_rotate"
  obj_type    = "AES"
  group_id    = dsm_group.group.id
  key_size    = 256
  key_ops     = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  rotate      = "DSM"
  // Name of the above security object
  rotate_from = "aes_sobject_example"
}
```

## Destroy a security object

```terraform
/*
Destruction can be done while update only. 
To destroy a security object, `destruct` parameter should be configured.
And make enabled as false, this is to avoid the differences while updating other resources.
*/
resource "dsm_sobject" "aes_sobject_example" {
  name            = "aes_sobject_example"
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
  description     = "aes sobject description"
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
    effective_at = "20241130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
  }
  destruct = "compromise" // other values: deactivate or destroy
  // Once compromised or destroyed, enabled will set to false. So, on terraform apply/plan, it ignores `enabled` parameter.
  // Once deactivated, original expiry date will set to the time of deactivation. So, on terraform apply/plan, it ignores `expiry_date` parameter.
  lifecycle {
    ignore_changes = [enabled, expiry_date]
  }
}
```