// Create a group
resource "dsm_group" "group" {
  name = "group"
  description = "group description"
}

// Import a secret
resource "dsm_secret" "secret" {
  name        = "secret"
  group_id    = dsm_group.group.id
  description = "test secret"
  enabled     = true
  state       = "Active"
  value       = "Rm9ydGFuaXg="
  expiry_date = "2025-02-02T17:04:05Z"
}

// Rotate a secret
resource "dsm_secret" "secret_rotate" {
  name        = "secret_rotate"
  group_id    = dsm_group.group.id
  description = "rotate secret"
  enabled     = true
  state       = "Active"
  value       = "cm90YXRlZm9ydGFuaXg="
  expiry_date = "2025-02-02T17:04:05Z"
  rotate      = true
  // Provide the secret name that needs to be rotated
  rotate_from = dsm_secret.secret.name
}