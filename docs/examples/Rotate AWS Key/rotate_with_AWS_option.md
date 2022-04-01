## Initially we need to create below mentioned 4 resources to create AWS BYOK key.

* Local Key Group

* Local Key

* AWS Group

* AWS Key

### Steps to create AWS BYOK key

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

## 1st Rotation of AWS key

Now we will rotate the AWS key `JCH-AWS-BYOK-KEY` created in above block.

### Steps to rotate AWS BYOK key

Append the below block to the existing Terraform file (above example) and apply the changes

* terraform apply

````
resource "dsm_aws_sobject" "test_sobject_aws_rotated" {
    name = "JCH-AWS-BYOK-KEY"
    group_id = dsm_aws_group.test_aws_group.id
    key = {
        kid = dsm_sobject.test_aes_key_rotated.id
    }
    rotate = "AWS"
    rotate_from = "JCH-AWS-BYOK-KEY"
}
````


## 2nd Rotation of AWS key

Now we will do 2nd rotation of the AWS key `JCH-AWS-BYOK-KEY` rotated in the above block.

### Steps to rotate AWS BYOK key

Append the below block to the existing Terraform file (above examples) and apply the changes

* terraform apply

````
resource "dsm_aws_sobject" "test_sobject_aws_rotated_2" {
    name = "JCH-AWS-BYOK-KEY"
    group_id = dsm_aws_group.test_aws_group.id
    key = {
        kid = dsm_sobject.test_aes_key_rotated_2.id
    }
    rotate = "AWS"
    rotate_from = "JCH-AWS-BYOK-KEY"
}
````