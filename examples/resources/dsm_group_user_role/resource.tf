# Adding a group role to a user
resource "dsm_group_user_role" "group_user_role" {
  group_name = "crypto_group"
  user_email = "test123@fortanix.com"
  role_name  = "GROUPAUDITOR"
}