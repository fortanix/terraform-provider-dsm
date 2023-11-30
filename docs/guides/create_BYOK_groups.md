## This examples will illustrate on how to create Cloud data control groups. We create AWS, Azure and GCP groups.

* AWS

```
// Create AWS group
resource "dsm_group" "aws_byok" {
    name = "aws_byok"
    description = "aws_byok"
    hmg = var.aws_data
}

// aws data to create a group inside dsm
variable "aws_data" {
  type        = any
  description = "The policy document. This is a JSON formatted string."
  default     = <<-EOF
    {
    "url": "kms.<region-names>.amazonaws.com",
    "tls": {
      "mode": "required",
      "validate_hostname": false,
      "ca": {
        "ca_set": "global_roots"
      }
    },
    "kind": "AWSKMS",
    "access_key": "<aws_access_key>",
    "secret_key": "<aws_secret_key>",
    "region": "<aws_region>",
    "service": "kms"
    }
  EOF
}
```
* Azure

```
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
      "url": "<key_vault_url>",
      "tls": {
        "mode": "required",
        "validate_hostname": false,
        "ca": {
          "ca_set": "global_roots"
        }
      },
      "kind": "AZUREKEYVAULT",
      "secret_key": "<aws_secret_key>",
      "tenant_id": "<azure_tenant_id>",
      "client_id": "<azure_client_id>",
      "subscription_id": "<azure_subscription_id>",
      "key_vault_type": "STANDARD"
    }
   EOF
}
```

* GCP

```
resource "dsm_group" "gcp_byok" {
  name = "gcp_byok"
  description = "gcp_byok"
  hmg = var.gcp_data
}

variable "gcp_data" {
  type        = any
  description = "The policy document. This is a JSON formatted string."
  default     = <<-EOF
    {
      "kind": "GCPKEYRING",
      "project_id": "<gcp-project-id>",
      "service_account_email": "<service-account-email-id>",
      "location": "<region-name>",
      "private_key": "<private-key-value>",
      "key_ring": "<key-ring-name>"
    }
   EOF
}
```