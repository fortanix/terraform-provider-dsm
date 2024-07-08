---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dsm_group Resource - terraform-provider-dsm"
subcategory: ""
description: |-
  Returns the Fortanix DSM group object from the cluster as a Resource.
---

# dsm_group (Resource)

Returns the Fortanix DSM group object from the cluster as a Resource.

## Example Usage

```terraform
resource "dsm_group" "group" {
  name = "group example"
  description = "group description"
  approval_policy = var.approval_policy
  hmg = var.azure_data
  key_undo_policy_window_time = 9000
}

variable "approval_policy" {
  type = any
  description = "The policy document. This is a JSON formatted string."
  default = <<-EOF
      {
        "protect_permissions": [
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
        ],
        "protect_crypto_operations": true,
        "quorum": {
          "n": 1,
          "members": [
            {
              "user": "54e489ca-f5aa-4e59-869e-281bbd37caa2"
            }
          ],
          "require_password": false,
          "require_2fa": false
        },
        "user": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "app": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
      }
      EOF
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The Fortanix DSM group object name.

### Optional

- `approval_policy` (String) The Fortanix DSM group object quorum approval policy definition as a JSON string
- `description` (String) The Fortanix DSM group object description
- `hmg` (String) The Fortanix DSM group object HMS/KMS definition as a JSON string
- `key_undo_policy_window_time` (Number) The Fortanix DSM group object key undo policy window time as an Integer(Number of seconds).

### Read-Only

- `acct_id` (String) Account ID from Fortanix DSM
- `creator` (Map of String) Creator of the group object from Fortanix DSM
- `group_id` (String) Group object ID from Fortanix DSM
- `hmg_id` (String) HSM/KMS ID from Fortanix/DSM
- `id` (String) The ID of this resource.