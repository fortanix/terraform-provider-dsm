variable "acct_id" {
  type    = string
  default = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}
# Add quorum policy to a Fortanix DSM account
# Example - 1

# When any one of the members' approval is required, assign `n` as 1 in the high level quorum.
# For example, members are users/apps that are configured in the Fortanix DSM account.
# The user/app value should be its UUID.
resource "dsm_acc_quorum_policy" "account_quorum_policy" {
  acct_id = var.acct_id
  approval_policy = jsonencode({
    policy = {
      quorum = {
        n = 1 # This defines that `n` member of approvals required.
        members = [
          {
            quorum = {
              n = 1, # This defines that `n` member of approvals required.
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                },
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
              require_2fa      = false
              require_password = true
            }
          },
          {
            quorum = {
              n = 1 # This defines that `n` member of approvals required.
              members = [
                {
                  app = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                },
                {
                  app = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
            }
          }
        ]
      }
    }
    manage_groups                  = false
    protect_authentication_methods = true
    protect_cryptographic_policy   = true
    protect_logging_config         = true
  })
}

# Add quorum policy to a Fortanix DSM account
# Example - 2
# When all the members of approval is required.
# In the example, members are users/apps that are configured in the Fortanix DSM account.
# The user/app value should be its UUID.
resource "dsm_acc_quorum_policy" "account_quorum_policy" {
  acct_id = var.acct_id
  approval_policy = jsonencode({
    policy = {
      quorum = {
        n = 2 # This defines that `n` member of approvals required.
        members = [
          {
            quorum = {
              n = 1 # This defines that `n` member of approvals required.
              members = [
                {
                  user = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
              require_2fa      = false
              require_password = true
            }
          },
          {
            quorum = {
              n = 1 # This defines that `n` member of approvals required.
              members = [
                {
                  app = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                }
              ]
            }
          }
        ]
      }
    }
    manage_groups                  = false
    protect_authentication_methods = true
    protect_cryptographic_policy   = true
    protect_logging_config         = true
  })
}