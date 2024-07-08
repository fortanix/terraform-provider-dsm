resource "dsm_group" "group" {
  name = "group example"
  description = "group description"
}

resource "dsm_secret" "secret" {
  name        = "Secret"
  group_id    = dsm_group.group.id
  description = "test secret"
  enabled     = true
  state       = "Active"
  value       = "Rm9ydGFuaXg="
  expiry_date = "20231130T183000Z"
}