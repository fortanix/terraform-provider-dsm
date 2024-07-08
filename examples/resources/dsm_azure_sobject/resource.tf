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

resource "dsm_azure_sobject" "sobject" {
  name            = "azure_sobject"
  group_id        = dsm_group.azure_byok.id
  description     = "key creation in akv"
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  enabled         = true
  expiry_date     = "20231130T183000Z"
  key             = {
    kid = "<dsm sobject key id>"
  }
  custom_metadata = {
    azure_key_state = "Enabled"
    azure-key-name = "key_inside_akv"
  }
  rotation_policy = {
    interval_days = 10
    effective_at = "20231130T183000Z"
    deactivate_rotated_key = true
  }
}