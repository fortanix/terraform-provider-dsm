# Step1: Read the existing DSM group that was already created.
resource "dsm_existing_group" "dsm_group" {
  name = "dsm_group"
}

# Step2: Update the group
# An example of a normal group where it did not configure hmg(HSM/KMS) while creating.
# hmg can not be added while updating the group.
resource "dsm_existing_group" "dsm_group" {
  name        = "dsm_group"
  description = "Update existing group"
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
      require_2fa      = false
    }
  })
}

# An example of a group where it configured the hmg(HSM/KMS) while creation.
resource "dsm_existing_group" "dsm_group" {
  name        = "dsm_group"
  description = "Update existing group"
  hmg = jsonencode({
    url = "https://sampleakv.vault.azure.net/"
    tls = {
      mode = "required"
      validate_hostname : false
      ca = {
        ca_set = "global_roots"
      }
    }
    kind            = "AZUREKEYVAULT"
    secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    tenant_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    client_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    subscription_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    key_vault_type  = "STANDARD"
  })
}