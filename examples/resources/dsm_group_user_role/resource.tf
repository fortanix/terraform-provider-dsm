resource "dsm_group_user_role" "sample_group_user_role" {
    group_name = "crypto_group_test"
    user_email = "test123@fortanix.com"
    role_name = "GROUPAUDITOR"
}