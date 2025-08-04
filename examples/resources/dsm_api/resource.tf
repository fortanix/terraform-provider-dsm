# Invoke a plugin
resource "dsm_api" "api" {
  method           = "POST"
  resource_type    = "plugin"
  resource_uuid    = "<plugin_id>"
  api_id_attribute = "kid" # Response of the API should have this attribute.
  payload = jsonencode({
    // payload request, eg:
    "key1" = "value1"
    "key2" = "value2"
  })
  recall = true
}
