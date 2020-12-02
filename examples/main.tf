terraform {
  required_providers {
    sdkms = {
      versions = ["0.1"]
      source = "fortanix.com/sdkms"
    }
  }
}

provider "sdkms" {
  endpoint = "https://sdkms.fortanix.com"
  username = "username"
  password = "password"
  acct_id  = "acct_id"
}

resource "sdkms_group" "group" {
  name     = "test-fyoo-group"
}

resource "sdkms_sobject" "sobject" {
  name     = "test-fyoo-sobject"
  key_size = 256
  obj_type = "AES" 
}