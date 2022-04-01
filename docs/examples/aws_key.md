## This example will create 4 resources.

* Local Key Group

* Local Key

* AWS Group

* AWS Key

## Steps to create AWS BYOK key

* terraform init 

* terraform apply

Note: Varibales used can also be defined in external vars.tf

```
terraform {
    required_providers {
        dsm = {
            version = "0.5.10"
            source = "fortanix/dsm"
        }
    }
}

provider "dsm" {    
    endpoint = var.endpoint
    username = var.username
    password = var.pass
    acct_id = var.acct_id
}

resource "dsm_group" "test_group" {
    name = "Local-BYOK-GROUP"
}

resource "dsm_sobject" "test_aes_key" {
    name            = "Local-AES-KEY"
    obj_type        = "AES"
    group_id        = dsm_group.test_group.id
    key_size        = 256
    key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
}

resource "dsm_aws_group" "test_aws_group" {
    name = "AWS-BYOK-GROUP"
    access_key = var.access_key
    secret_key = var.secret_key
}

resource "dsm_aws_sobject" "test_sobject_aws" {
    name = "JCH-AWS-BYOK-KEY"
    group_id = dsm_aws_group.test_aws_group.id
    key = {
        kid = dsm_sobject.test_aes_key.id
    }
    custom_metadata = {
        aws-aliases = "JCH-AWS-BYOK-KEY"
        aws-policy = <custom-aws-policy>
    }
}
```