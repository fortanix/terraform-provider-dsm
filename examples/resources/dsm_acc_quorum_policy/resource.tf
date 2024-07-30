variable "acct_id" {
  type    = string
  default = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}
// Add quorum policy to a Fortanix DSM account
// Example - 1
/*
When any one of the members approval is required.
In the example, members are users that are configured in the Fortanix DSM account.
Apps can also be the members. e.g. { "app": "<app_uuid>" }
*/
resource "dsm_acc_quorum_policy" "account_quorum_policy" {
  acct_id = var.acct_id
  approval_policy = jsonencode({
    policy = {
      quorum = {
        n = 1
        members = [
          {
            quorum = {
              n = 1,
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
              require_2fa = false
              require_password = true
            }
          },
          {
            quorum = {
              n = 1
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
            }
          }
        ]
      }
    }
    manage_groups = false
    protect_authentication_methods = true
    protect_cryptographic_policy = true
    protect_logging_config = true
  })
}

// Add quorum policy to a Fortanix DSM account
// Example - 2
/*When all the members of approval is required.
In the example, members are users that are configured in the Fortanix DSM account.
Apps can also be the members. Apps can also be the members. e.g. { "app": "<app_uuid>" }
*/
resource "dsm_acc_quorum_policy" "account_quorum_policy" {
  acct_id         = var.acct_id
  approval_policy = jsonencode({
    policy = {
      quorum = {
        n = 2
        members = [
          {
            quorum = {
              n = 1
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
              require_2fa = false
              require_password = true
            }
          },
          {
            quorum = {
              n = 1
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
            }
          }
        ]
      }
    }
    manage_groups = false
    protect_authentication_methods = true
    protect_cryptographic_policy = true
    protect_logging_config = true
  })
}