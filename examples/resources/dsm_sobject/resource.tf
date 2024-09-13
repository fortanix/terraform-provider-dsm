# Create a group
resource "dsm_group" "group" {
  name        = "group"
  description = "group description"
}
# Create a security object in the above group
resource "dsm_sobject" "sobject" {
  name     = "sobject"
  obj_type = "AES"
  group_id = dsm_group.group.id
  key_size = 256
  key_ops = [
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
  enabled     = true
  expiry_date = "2025-02-02T17:04:05Z"
  description = "sobject description"
  custom_metadata = {
    key1 = "value1"
  }
  allowed_key_justifications_policy = [
    "CUSTOMER_INITIATED_SUPPORT",
    "CUSTOMER_INITIATED_ACCESS"
  ]
  allowed_missing_justifications = true
  rotation_policy = {
    interval_days          = 20
    effective_at           = "20231130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys     = "all_external"
  }
}

## copy a security object

# Copy above security object.
# When copying a security object obj_type, key_size, allowed_key_justifications_policy, allowed_key_justifications_policy, -
# allowed_missing_justifications, lms or bls should not be configured.

resource "dsm_sobject" "sobject_copy" {
  name     = "sobject_copy"
  group_id = dsm_group.group.id
  key = {
    kid = dsm_sobject.sobject.id
  }
  key_ops = [
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
  enabled     = true
  expiry_date = "2026-02-02T17:04:05Z"
  description = "sobject description copy"
  custom_metadata = {
    key1 = "value1"
  }
  rotation_policy = {
    interval_days          = 20
    effective_at           = "20231130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys     = "all_external"
  }
}

# rotate a security object
resource "dsm_sobject" "sobject_rotate" {
  name     = "sobject_rotate"
  obj_type = "AES"
  group_id = dsm_group.group.id
  key_size = 256
  key_ops  = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  rotate   = "DSM"
  # Name of the above security object
  rotate_from = "sobject"
}

## import a security object

# This is an example of importing a certificate
resource "dsm_sobject" "certificate" {
  name        = "certificate"
  obj_type    = "CERTIFICATE"
  group_id    = dsm_group.group.id
  value       = "XXXXXXXXXXXX<CERTIFICATE_value_in_a_string_format>XXXXXXXXXXXXXX"
  expiry_date = "2025-02-02T17:04:05Z"
  enabled     = true
  key_ops = [
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