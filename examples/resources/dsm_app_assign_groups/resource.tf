# create a new group1
resource "dsm_group" "group1" {
  name = "group1"
}

# create a new group2
resource "dsm_group" "group2" {
  name = "group2"
}

# Assign group1 and group2 to an existing App using app name
resource "dsm_app_assign_groups" "dsm_app_assign_groups" {
  app_name = "ExistingApp"
  groups   = [dsm_group.group1.id, dsm_group.group2.id]
}

# create a new group3
resource "dsm_group" "group3" {
  name = "group3"
}

# create a new group3
resource "dsm_group" "group4" {
  name = "group4"
}

# Assign group1 and group2 to an existing App using app id
resource "dsm_app_assign_groups" "dsm_app_assign_groups" {
  app_id = "12a3bd2c-0e23-46e9-bcf8-57526602f629"
  groups = [dsm_group.group3.id, dsm_group.group4.id]
}
