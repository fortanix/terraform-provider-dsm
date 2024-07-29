// Create a normal group
resource "dsm_group" "normal_group" {
  name = "group_test"
}

// Create GCP group
resource "dsm_group" "gcp_group" {
  name = "gcp_group"
  hmg = jsonencode({
    kind         = "GCPKEYRING"
    key_ring       = "key_ring_name"
    project_id      = "gcp_project_id"
    service_account_email = "test@test.iam.gserviceaccount.com"
    location       = "us-east1"
    private_key      = "<Private component of the service account key pair that can be obtained from the GCP cloud console. It is used to authenticate the requests made by DSM to the GCP cloud. This should be base64 encoded private key.>"
  })
}

// Create an AES key in normal group
resource "dsm_sobject" "sobject" {
  name     = "aes256"
  key_size = 256
  group_id = dsm_group.normal_group.id
  obj_type    = "AES"
}

// Copy a key to GCP key ring using the above DSM security object
resource "dsm_gcp_sobject" "gcp_sobject" {
  name     = "gcp_sobject"
  group_id = dsm_group.gcp_group.id
  key = {
    kid = dsm_sobject.sobject.id
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
  obj_type    = "AES"
  key_size    = 256
  key_ops     = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE", "EXPORT"]
  enabled     = true
  expiry_date = "2025-02-02T17:04:05Z"
}
