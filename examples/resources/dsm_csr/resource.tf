# To use this resource, it is required to create a custom plugin in DSM first.
# Copy the plugin from https:#github.com/fortanix/terraform-provider-dsm/blob/main/plugins/Terraform-Plugin-CSR.lua
# Create the custom plugin in DSM 
# Plugin title: "Terraform Plugin - CSR"


# Create an RSA key pair that will be used to generate the CSR
resource "dsm_sobject" "sobject" {
  name     = "sobject-rsa"
  obj_type = "RSA"
  group_id = "<group ID>" # make sure that the group can be accessed by your plugin "Terraform Plugin - CSR".
  key_size = 2048
}

# Generating the CSR
resource "dsm_csr" "csr" {
  kid      = dsm_sobject.sobject.id
  cn       = "example-common-name"
  ou       = "example-organizational-unit"
  o        = "example-organization"
  l        = "example-location"
  c        = "example-country"
  st       = "example-state"
  e        = "example@example.com"
  email    = ["alt-email@example.com"]
  dnsnames = ["example.com", "www.example.com"]
  ips      = ["192.168.1.1", "10.0.0.1"]
}