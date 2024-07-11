resource "dsm_group" "group" {
  name = "group example"
  description = "group description"
}

resource "dsm_sobject" "sobject" {
  name            = "sobject_example"
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
  expiry_date     = "20231130T183000Z"
  description     = "sobject description"
  custom_metadata = {
    key1 = "value1"
  }
  allowed_key_justifications_policy = [
    "CUSTOMER_INITIATED_SUPPORT",
    "CUSTOMER_INITIATED_ACCESS",
    "CUSTOMER_AUTHORIZED_WORKFLOW_SERVICING",
    "GOOGLE_INITIATED_SERVICE",
    "GOOGLE_INITIATED_REVIEW",
    "GOOGLE_INITIATED_SYSTEM_OPERATION",
    "THIRD_PARTY_DATA_REQUEST",
    "REASON_UNSPECIFIED",
    "REASON_NOT_EXPECTED",
    "MODIFIED_CUSTOMER_INITIATED_ACCESS",
    "MODIFIED_GOOGLE_INITIATED_SYSTEM_OPERATION",
    "GOOGLE_RESPONSE_TO_PRODUCTION_ALERT"
  ]
  allowed_missing_justifications = true
  rotation_policy = {
    interval_days = 20
    effective_at = "20231130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys = "all_external"
}
}