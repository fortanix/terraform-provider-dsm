## This examples will illustrate on how to create Cloud data control groups. We create AWS, Azure and GCP groups.

* AWS

```
// Create AWS group
resource "dsm_group" "aws_group" {
  name = "aws_group"
  description = "AWS group"
  hmg = jsonencode(
    {
      url = "kms.us-east-1.amazonaws.com"
      tls = {
        mode = "required"
        validate_hostname: false,
        ca = {
          ca_set = "global_roots"
        }
      }
      kind = "AWSKMS"
      access_key = "XXXXXXXXXXXXXXXXXXXX"
      secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
      region = "us-east-1"
      service = "kms"
    })
}
```
* Azure

```
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

```

* GCP

```
// Create GCP group
resource "dsm_group" "gcp_group" {
  name = "gcp_group"
  hmg = jsonencode({
    kind         = "GCPKEYRING"
    key_ring       = "key_ring_name"
    project_id      = "gcp_project_id"
    service_account_email = "test@test.iam.gserviceaccount.com"
    location       = "us-east1"
    private_key      = "<Private component of the service account key pair that can be obtained from the GCP cloud console. It is used to authenticate the requests made by DSM to the GCP cloud. This should be base64 encoded private key.>"
  })
}
```