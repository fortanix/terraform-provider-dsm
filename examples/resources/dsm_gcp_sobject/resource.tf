resource "dsm_gcp_sobject" "sample_gcp_sobject" {
  name     = "test-gcp-sobject"
  group_id = "311915f2-7cdb-4ea9-ac15-83818f04dc39"
  key = {
    kid = "f6be4755-7912-4546-9e94-27851f2ddcd7"
  }
  custom_metadata = {
    gcp-key-id = "name-of-the-key-in-gcp"
  }
  rotation_policy = {
    interval_days          = 20
    effective_at           = "20231130T183000Z"
    deactivate_rotated_key = true
    rotate_copied_keys     = "all_external"
  }
  obj_type = "AES"
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
}