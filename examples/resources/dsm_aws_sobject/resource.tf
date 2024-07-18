// Create a normal group
resource "dsm_group" "normal_group" {
  name = "normal_group"
}

// Create AWS group
resource "dsm_group" "aws_group" {
  name = "aws_group"
  description = "AWS group change"
  hmg = var.aws_data
}

// aws credential data to create a group inside dsm
variable "aws_data" {
  type        = any
  description = "The policy document. This is a JSON formatted string."
  default     = <<-EOF
    {
    "url": "kms.us-east-1.amazonaws.com",
    "tls": {
      "mode": "required",
      "validate_hostname": false,
      "ca": {
        "ca_set": "global_roots"
      }
    },
    "kind": "AWSKMS",
    "access_key": "XXXXXXXXXXXXXXXXXXXX",
    "secret_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "region": "us-east-1",
    "service": "kms"
    }
  EOF
}

// Create an AES key inside DSM
resource "dsm_sobject" "aes_sobject" {
  name            = "aes_sobject"
  obj_type        = "AES"
  group_id        = dsm_group.normal_group.id
  key_size        = 256
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
}

// AWS sobject creation(Copies the key from DSM)
resource "dsm_aws_sobject" "aws_sobject" {
  name = "aws_sobject"
  group_id = dsm_group.aws_group.id
  description = "AWS sobject"
  enabled = true
  expiry_date     = "2025-02-02T17:04:05Z"
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  key = {
    kid = dsm_sobject.aes_sobject.id
  }
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
    aws-policy = "{\"Version\":\"2012-10-17\",\"Id\":\"key-default-1\",\"Statement\":[{\"Sid\":\"EnableIAMUserPermissions\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"arn:aws:iam::XXXXXXXXXXXX:root\"},\"Action\":\"kms:*\",\"Resource\":\"*\"}]}"
  }
}

// Note: For rotation of a key please refer Guides/rotate_with_AWS_option.
