// Create a normal group
resource "dsm_group" "group" {
  name = "group"
}

/*
Create a group with multiple parameters.
The following resource group is an example of an external KMS group of Azure key vault
and an approval policy.

For more examples of external KMS groups please refer Guides/create_BYOK_groups
*/
resource "dsm_group" "group" {
  name = "group"
  description = "group description"
  approval_policy = jsonencode({
    protect_permissions = [
      "ROTATE_SOBJECTS",
      "REVOKE_SOBJECTS",
      "REVERT_SOBJECTS",
      "DELETE_KEY_MATERIAL",
      "DELETE_SOBJECTS",
      "DESTROY_SOBJECTS",
      "MOVE_SOBJECTS",
      "CREATE_SOBJECTS",
      "UPDATE_SOBJECTS_PROFILE",
      "UPDATE_SOBJECTS_ENABLED_STATE",
      "UPDATE_SOBJECT_POLICIES",
      "ACTIVATE_SOBJECTS",
      "UPDATE_KEY_OPS"
    ]
    protect_crypto_operations = true
    quorum = {
      n = 1
      members = [
        {
          user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
        }
      ]
      require_password = false
      require_2fa = false
    }
  })
  hmg = jsonencode({
    url = "https://sampleakv.vault.azure.net/"
    tls = {
      mode = "required"
      validate_hostname : false
      ca = {
        ca_set = "global_roots"
      }
    }
    kind = "AZUREKEYVAULT"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    tenant_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    client_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    subscription_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    key_vault_type = "STANDARD"
  })
  key_undo_policy_window_time = 9000
}