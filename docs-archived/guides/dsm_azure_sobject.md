***Rotation of dsm_azure_sobject***

**Rotate with DSM Option**

```terraform
// Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}
// Create an Azure BYOK group
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

// Create an RSA security object in normal group
resource "dsm_sobject" "rsa_key_dsm" {
  name     = "rsa_key_dsm"
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
}

// Copy above RSA security object to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure" {

  name            = "rsa_key_azure"
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
}

// 1st Rotation of azure security object with DSM option

// Rotate RSA security object
resource "dsm_sobject" "rsa_key_dsm_rotate1" {
  name     = dsm_sobject.rsa_key_dsm.name
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops  = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
  rotate   = "DSM"
  rotate_from = dsm_sobject.rsa_key_dsm.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}

// Copy above RSA key to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure_rotate1" {
  name            = dsm_azure_sobject.rsa_key_azure.name
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm_rotate1.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate   = "DSM"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}


// 2nd Rotation of azure security object with DSM option

// Rotate RSA security object
resource "dsm_sobject" "rsa_key_dsm_rotate2" {
  name     = dsm_sobject.rsa_key_dsm.name
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops  = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
  rotate   = "DSM"
  rotate_from = dsm_sobject.rsa_key_dsm.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}

// Copy above RSA key to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure_rotate2" {
  name            = dsm_azure_sobject.rsa_key_azure.name
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm_rotate2.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate   = "DSM"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}
```

**Rotate with AZURE Option**

```terraform
// Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}
// Create an Azure BYOK group
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

// Create an RSA security object in normal group
resource "dsm_sobject" "rsa_key_dsm" {
  name     = "rsa_key_dsm"
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
}

// Copy above RSA security object to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure" {
  name            = "rsa_key_azure"
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
}

// 1st Rotation of azure security object with AZURE option

// Rotate RSA security object
resource "dsm_sobject" "rsa_key_dsm_rotate1" {
  name     = dsm_sobject.rsa_key_dsm.name
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops  = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
  rotate   = "AZURE"
  rotate_from = dsm_sobject.rsa_key_dsm.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}

// Copy above RSA key to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure_rotate1" {
  name            = dsm_azure_sobject.rsa_key_azure.name
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm_rotate1.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate   = "AZURE"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}


// 2nd Rotation of azure security object with AZURE option

// Rotate RSA security object
resource "dsm_sobject" "rsa_key_dsm_rotate2" {
  name     = dsm_sobject.rsa_key_dsm.name
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops  = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
  rotate   = "AZURE"
  rotate_from = dsm_sobject.rsa_key_dsm.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}

// Copy above RSA key to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure_rotate2" {
  name            = dsm_azure_sobject.rsa_key_azure.name
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm_rotate2.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate   = "AZURE"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name
  lifecycle {
    ignore_changes = [rotate, rotate_from]
  }
}
```

***Schedule deletion of AWS security object***

```terraform
// Schedule an DSM Azure security object to delete
/*
Enable schedule_deletion as true.
This can be enabled during both creation and updation. 
*/
resource "dsm_azure_sobject" "rsa_key_azure" {
  name            = "rsa_key_azure"
  group_id        = dsm_group.azure_group.id
  key             = {
    kid =  dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-state =  "Enabled"
    azure-key-name = "rsa-key-azure"
  }
  key_ops         = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  schedule_deletion = true
}

```


