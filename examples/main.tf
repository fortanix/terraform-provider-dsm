terraform {
  required_providers {
    sdkms = {
      version = "0.1.2"
      source = "fortanix.com/fyoo/sdkms"
    }
  }
}

data "external" "aws-sts-generator" {
  //program = ["bash", "-c" "aws", "sts", "assume-role", "--role-arn", "arn:aws:iam::513076507034:role/aws-kms-power-user", "--role-session-name", "terraform-access", "--output", "json", "|", "jq", ".Credentials"]
  program = ["bash", "-c", "aws sts assume-role --role-arn arn:aws:iam::513076507034:role/aws-kms-power-user --role-session-name terraform-test --output json | jq .Credentials"]
}

output "something" {
  value = data.external.aws-sts-generator.result.AccessKeyId
}

provider "sdkms" {
  endpoint = "https://sdkms.fortanix.com"
  username = ""
  password = ""
  acct_id  = ""
}

resource "sdkms_group" "group" {
  name     = "test-fyoo-group"
}

resource "sdkms_aws_group" "awsgroup" {
  name     = "test-fyoo-awsgroup"
  url      = "kms.us-east-2.amazonaws.com"
  access_key = data.external.aws-sts-generator.result.AccessKeyId
  secret_key = data.external.aws-sts-generator.result.SecretAccessKey
}

resource "sdkms_app" "app" {
  name     = "test-fyoo-app"
  default_group = sdkms_group.group.id
}

resource "sdkms_sobject" "sobject" {
  name     = "test-fyoo-sobject"
  group_id = sdkms_group.group.id
  key_size = 256
  obj_type = "AES" 
}

resource "sdkms_sobject" "tokenis" {
  name     = "test-fyoo-tokenis"
  group_id = sdkms_group.group.id
  key_size = 256
  obj_type = "AES"
  key_ops  =  [
    "ENCRYPT",
    "DECRYPT",
    "APPMANAGEABLE" ]
  fpe_radix = 10 
}
