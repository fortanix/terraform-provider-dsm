---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dsm_azure_sobject Resource - terraform-provider-dsm"
subcategory: ""
description: |-
  Creates a new security object in Azure key vault. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to Azure KV as a Customer Managed Key (CMK).
---

# dsm_azure_sobject (Resource)

Creates a new security object in Azure key vault. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to Azure KV as a Customer Managed Key (CMK).

## Example Usage

```terraform
// Create Azure group
resource "dsm_group" "azure_group" {
  name = "azure_group"
  description = "azure_group"
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

/* Copy a key to azure key vault using the above DSM security object.
By default it creates a key as a software protected key.
*/
resource "dsm_azure_sobject" "azure_sobject" {
  name            = "azure_sobject"
  group_id        = dsm_group.azure_group.id
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

/* Copy a key to azure key vault using the above DSM security object.
It is an example of hardware protected key in PREMIUM key vault.
*/
resource "dsm_azure_sobject" "sobject" {
  name            = "azure_sobject"
  group_id        = dsm_group.azure_group.id
  description     = "key creation in akv"
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT", "APPMANAGEABLE"]
  enabled         = true
  expiry_date     = "2025-02-02T17:04:05Z"
  key             = {
    kid = dsm_sobject.dsm_sobject.id
  }
  custom_metadata = {
    azure-key-name = "key_inside_akv"
    azure-key-type = "hardware"
  }
  rotation_policy = {
    interval_days = 10
    effective_at = "20231130T183000Z"
    deactivate_rotated_key = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `custom_metadata` (Map of String) Azure CMK level metadata information.
   * `azure-key-name`: Key name within Azure KV.
   * **Note:** By default dsm_azure_sobject creates the key as a software protected key. For a hardware protected key use the below parameter.
   * `azure-key-type`: Type of a key. It can be used in `PREMIUM` key vault. Value is hardware.
- `group_id` (String) The Azure group ID in Fortanix DSM into which the key will be generated.
- `key` (Map of String) A local security object imported to Fortanix DSM(BYOK) and copied to Azure KV.
- `name` (String) The security object name.

### Optional

- `description` (String) The security object description.
- `enabled` (Boolean) Whether the security object will be Enabled or Disabled. The values are true/false.
- `expiry_date` (String) The security object expiry date in RFC format.
- `key_ops` (List of String) The security object operations permitted.

| obj_type | key_size/curve | key_ops |
| -------- | -------- |-------- |
| `RSA` | 2048, 3072, 4096 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |
| `EC` | NistP256, NistP384, NistP521,SecP256K1 | APPMANAGEABLE, SIGN, VERIFY, AGREEKEY, EXPORT
- `key_size` (Number) The size of the security object.
- `obj_type` (String) The type of security object.
- `rotation_policy` (Map of String) Policy to rotate a Security Object, configure the below parameters.
   * `interval_days`: Rotate the key for every given number of days.
   * `interval_months`: Rotate the key for every given number of months.
   * `effective_at`: Start of the rotation policy time.
   * `deactivate_rotated_key`: Deactivate original key after rotation true/false.
   * **Note:** Either interval_days or interval_months should be given, but not both.
- `state` (String) The key states of the Azure KV key. The values are Created, Deleted, Purged.

### Read-Only

- `acct_id` (String) The account ID from Fortanix DSM.
- `creator` (Map of String) The creator of the security object from Fortanix DSM.
   * `user`: If the security object was created by a user, the computed value will be the matching user id.
   * `app`: If the security object was created by a app, the computed value will be the matching app id.
- `id` (String) The ID of this resource.
- `kid` (String) The security object ID from Fortanix DSM.
- `links` (Map of String) Link between local security object and Azure KV security object.
