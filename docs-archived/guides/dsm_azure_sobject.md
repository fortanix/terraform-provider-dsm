# Rotation of dsm_azure_sobject

## Rotate with DSM Option

```terraform
# Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}
# Create an Azure BYOK group
resource "dsm_group" "azure_group" {
  name        = "azure_group"
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
    kind            = "AZUREKEYVAULT"
    secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    tenant_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    client_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    subscription_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    key_vault_type  = "STANDARD"
  })
}

# Create an RSA security object in normal group
resource "dsm_sobject" "rsa_key_dsm" {
  name     = "rsa_key_dsm"
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops  = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
}

# Copy above RSA security object to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure" {
  name     = "rsa_key_azure"
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure"
  }
  key_ops = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
}

## 1st Rotation of azure security object with DSM option

# Just copy the above dsm_azure_sobject.rsa_key_azure, add rotate and rotate_from attributes and change the resource name
resource "dsm_azure_sobject" "rsa_key_azure_rotate1" {
  name     = dsm_azure_sobject.rsa_key_azure.name # Name should be the same as the key to be rotated.
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure" # Should be the same azure-key-name as above
  }
  key_ops     = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate      = "DSM"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name # Name of the key from where it is being rotated.
}


## 2nd Rotation of azure security object with DSM option

# Just copy the above dsm_azure_sobject.rsa_key_azure_rotate1, add rotate and rotate_from attributes and change the resource name
resource "dsm_azure_sobject" "rsa_key_azure_rotate2" {
  name     = dsm_azure_sobject.rsa_key_azure_rotate1.name # Name should be the same as the key name to be rotated.
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure" # Should be the same azure-key-name as above 
  }
  key_ops     = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate      = "DSM"
  rotate_from = dsm_azure_sobject.rsa_key_azure_rotate1.name # Name of the key from where it is being rotated.
}

```

## Rotate with AZURE Option

```terraform
# Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}
# Create an Azure BYOK group
resource "dsm_group" "azure_group" {
  name        = "azure_group"
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
    kind            = "AZUREKEYVAULT"
    secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    tenant_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    client_id       = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    subscription_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
    key_vault_type  = "STANDARD"
  })
}

# Create a RSA security object in normal group
resource "dsm_sobject" "rsa_key_dsm" {
  name     = "rsa_key_dsm"
  group_id = dsm_group.normal_group.id
  key_size = 2048
  key_ops = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "SIGN", "VERIFY", "EXPORT"]
  obj_type = "RSA"
}

# Copy above RSA security object to azure key vault
resource "dsm_azure_sobject" "rsa_key_azure" {
  name     = "rsa_key_azure"
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure"
  }
  key_ops = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
}

## 1st Rotation of azure security object with AZURE option

# Just copy the above dsm_azure_sobject.rsa_key_azure, add rotate and rotate_from attributes and change the resource name
resource "dsm_azure_sobject" "rsa_key_azure_rotate1" {
  name     = dsm_azure_sobject.rsa_key_azure.name # Name should be the same as the key name to be rotated.
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure" # Should be the same azure-key-name as above
  }
  key_ops     = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate      = "AZURE"
  rotate_from = dsm_azure_sobject.rsa_key_azure.name # Name of the key from where it is being rotated.
}

## 2nd Rotation of azure security object with AZURE option

# Just copy the above dsm_azure_sobject.rsa_key_azure_rotate1, add rotate and rotate_from attributes and change the resource name
resource "dsm_azure_sobject" "rsa_key_azure_rotate2" {
  name     = dsm_azure_sobject.rsa_key_azure_rotate1.name # Name should be the same as the key name to be rotated.
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure" # Should be the same azure-key-name as above
  }
  key_ops     = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  rotate      = "AZURE"
  rotate_from = dsm_azure_sobject.rsa_key_azure_rotate1.name # Name of the key from where it is being rotated.
}
```

## Soft deletion and Purge key of an Azure security object

```terraform
## Soft deletion of dsm_azure_sobject

# Enable soft_deletion as true.
# This can be enabled only during update.
resource "dsm_azure_sobject" "rsa_key_azure" {
  name     = "rsa_key_azure"
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure"
  }
  key_ops       = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  soft_deletion = true
}

## Purging a dsm_azure_sobject.

# Enable purge_deleted_key as true.
# This can be enabled only during update and can be done only after soft_deletion.
resource "dsm_azure_sobject" "rsa_key_azure" {
  name     = "rsa_key_azure"
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure"
  }
  key_ops           = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  purge_deleted_key = true
}

## Soft deletion and Purging a key in a single request.

# First it does the soft deletion and then purging the key.
# These can be enabled only during update.
resource "dsm_azure_sobject" "rsa_key_azure" {
  name     = "rsa_key_azure"
  group_id = dsm_group.azure_group.id
  key = {
    kid = dsm_sobject.rsa_key_dsm.id
  }
  custom_metadata = {
    azure-key-name  = "rsa-key-azure"
  }
  key_ops           = ["SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "EXPORT", "APPMANAGEABLE", "HIGHVOLUME"]
  soft_deletion     = true
  purge_deleted_key = true
}
```


