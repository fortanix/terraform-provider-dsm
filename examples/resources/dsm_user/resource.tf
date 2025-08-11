# Create a User with the role ACCOUNTADMINISTRATOR
resource "dsm_user" "user1" {
  user_email  = "user1@test.com"
  description = "user1"
  first_name  = "user1"
  last_name   = "test"
  role        = "ACCOUNTADMINISTRATOR"
}

# Create a User with the role ACCOUNTAUDITOR
resource "dsm_user" "user2" {
  user_email  = "user2@test.com"
  description = "user2"
  first_name  = "user2"
  last_name   = "test"
  role        = "ACCOUNTAUDITOR"
}

# Create a User with the role ACCOUNTMEMBER

## Create groups
resource "dsm_group" "group1" {
  name = "group1"
}

resource "dsm_group" "group2" {
  name = "group2"
}

resource "dsm_group" "group3" {
  name = "group3"
}

resource "dsm_group" "group4" {
  name = "group4"
}


resource "dsm_user" "user3" {
  user_email  = "user3@test.com"
  description = "user3"
  first_name  = "user3"
  last_name   = "test"
  role        = "ACCOUNTMEMBER"
  # Add specific groups and their roles
  groups = jsonencode({
    "${dsm_group.group1.id}" = ["GROUPADMINISTRATOR"]
    "${dsm_group.group2.id}" = ["GROUPAUDITOR"]
    "${dsm_group.group3.id}" = ["GROUPAUDITOR"]
    "${dsm_group.group4.id}" = ["GROUPADMINISTRATOR"]

  })
}