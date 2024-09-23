# Creation of azure group
resource "dsm_azure_group" "dsm_azure_group" {
  name            = "dsm_azure_group"
  description     = "Azure group"
  url             = "https:#testfortanixterraform.vault.azure.net/"
  tenant_id       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  client_id       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  subscription_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  key_vault_type  = "STANDARD"
}