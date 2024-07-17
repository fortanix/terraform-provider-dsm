resource "dsm_csr" "example_csr" {
  kid = "e5c95efc-e08d-4274-936e-da8e3d88286d"
  cn = "example-common-name"
  ou = "example-organizational-unit"
  o = "example-organization"
  l = "example-location"
  c = "example-country"
  st = "example-state"
  e = "example@example.com"
  email = [ "alt-email@example.com" ]
  dnsnames = [ "example.com", "www.example.com" ]
  ips = [ "192.168.1.1", "10.0.0.1" ]
}