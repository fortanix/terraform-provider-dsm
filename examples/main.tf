terraform {
  required_providers {
    sdkms = {
      version = "0.1.2"
      source = "fortanix.com/fyoo/sdkms"
    }
  }
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
