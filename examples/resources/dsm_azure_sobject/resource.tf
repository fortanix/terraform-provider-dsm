// Create Azure group
resource "dsm_group" "azure_byok" {
  name = "azure_byok"
  description = "azure_byok"
  hmg = var.azure_data
}

// azure data to create a group inside dsm
variable "azure_data" {
  type        = any
  description = "The policy document. This is a JSON formatted string."
  default     = <<-EOF
    {
      "url": "https://sampleakv.vault.azure.net/",
      "tls": {
        "mode": "required",
        "validate_hostname": false,
        "ca": {
          "ca_set": "global_roots"
        }
      },
      "kind": "AZUREKEYVAULT",
      "secret_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
      "tenant_id": "0XXXXXXX-YYYY-HHHH-GGGG-123456789123",
      "client_id": "0XXXXXXX-YYYY-HHHH-GGGG-123456789123",
      "subscription_id": "0XXXXXXX-YYYY-HHHH-GGGG-123456789123",
      "key_vault_type": "STANDARD"
    }
   EOF
}

// Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}

// Create a RSA key in normal group
resource "dsm_sobject" "dsm_sobject" {
  name     = "dsm_sobject"
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT", "APPMANAGEABLE"]
  obj_type = "RSA"
}

// Copy a key to azure key vault using the above DSM security object
resource "dsm_azure_sobject" "sobject" {
  name            = "azure_sobject"
  group_id        = dsm_group.azure_byok.id
  description     = "key creation in akv"
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT", "APPMANAGEABLE"]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  key             = {
    kid = dsm_sobject.dsm_sobject.id
  }
  custom_metadata = {
    azure-key-name = "key_inside_akv"
  }
  rotation_policy = {
    interval_days = 10
    effective_at = "20231130T183000Z"
    deactivate_rotated_key = true
  }
}