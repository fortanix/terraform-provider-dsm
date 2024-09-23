# Creation of a groups
resource "dsm_group" "group1" {
  name = "group1"
}

resource "dsm_group" "group2" {
  name = "group2"
}

resource "dsm_group" "group3" {
  name = "group3"
}

# Read the lua plugin from a file
data "local_file" "plugin_code" {
  filename = "path/of/a/lua_plugin"
}

# Create a plugin by reading a file
resource "dsm_plugin" "dsm_plugin" {
  name          = "dsm_plugin"
  description   = "DSM Plugin"
  default_group = dsm_group.group1.id
  groups        = [dsm_group.group1.id, dsm_group.group2.id, dsm_group.group3.id]
  plugin_type   = "STANDARD"
  language      = "LUA"
  code          = data.local_file.plugin_code.content
}