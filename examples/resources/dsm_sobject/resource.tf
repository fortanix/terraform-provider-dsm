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