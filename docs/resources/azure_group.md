---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dsm_azure_group Resource - terraform-provider-dsm"
subcategory: ""
description: |-
  Creates a Fortanix DSM group mapped to Azure Key Vault in the cluster as a resource. This group acts as a container for security objects. The returned resource object contains the UUID of the group for further references.
---

# dsm_azure_group (Resource)

Creates a Fortanix DSM group mapped to Azure Key Vault in the cluster as a resource. This group acts as a container for security objects. The returned resource object contains the UUID of the group for further references.

## Example Usage

```terraform
# Creation of azure group
resource "dsm_azure_group" "dsm_azure_group" {
  name            = "dsm_azure_group"
  description     = "Azure group"
  url             = "https://testfortanixterraform.vault.azure.net/"
  tenant_id       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  client_id       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  subscription_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  key_vault_type  = "STANDARD"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) The Azure registered application id (username).
- `name` (String) The Azure KV group object name in Fortanix DSM.
- `secret_key` (String, Sensitive) A secret string that a registered application in Azure uses to prove its identity (application password).
- `subscription_id` (String) The ID of the Azure AD subscription.
- `tenant_id` (String) The tenant/directory id of the Azure subscription.
- `url` (String) The URL of the object in an Azure KV that uniquely identifies the object.

### Optional

- `description` (String) Description of the Azure KV Fortanix DSM group.
- `key_vault_type` (String) The type of key vault. The default value is `Standard`. Values are Standard/Premium.

### Read-Only

- `acct_id` (String) The Account ID from Fortanix DSM.
- `creator` (Map of String) The creator of the group from Fortanix DSM.
   * `user`: If the group was created by a user, the computed value will be the matching user id.
   * `app`: If the group was created by a app, the computed value will be the matching app id.
- `group_id` (String) The Azure KV group object ID from Fortanix DSM.
- `id` (String) The ID of this resource.
