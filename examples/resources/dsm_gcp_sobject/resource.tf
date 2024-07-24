// Create GCP group
resource "dsm_group" "gcp_cdc" {
  name = "group_gcp"
  hmg  = var.hmg
}

// GCP data to create a group inside DSM
variable "hmg" {
  default = <<EOF
  {
    "kind": "GcpKeyRing",
    "service_account_email": "test@test.iam.gserviceaccount.com",
    "project_id": "fortanix",
    "location": "us-east1",
    "key_ring": "key_ring_name",
    "private_key": "<Private component of the service account key pair that can be obtained from the GCP cloud console. It is used to authenticate the requests made by DSM to the GCP cloud>"
  }
  EOF
}

// Create a normal group
resource "dsm_group" "normal_group" {
  name = "group_test"
}

// Create an AES key in normal group
resource "dsm_sobject" "sobject" {
  name     = "aes256"
  key_size = 256
}

// Copy a key to GCP key ring using the above DSM security object
resource "dsm_gcp_sobject" "sample_gcp_sobject" {
  name     = "test-gcp-sobject"
  group_id = dsm_group.gcp_cdc.id
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
