// Add quorum policy to a Fortanix DSM account
// Example - 1
resource "dsm_acc_quorum_policy" "account_quorum_policy" {
  acct_id = var.acct_id
  approval_policy = var.account_quorum_policy1
}

variable "acct_id" {
  type    = string
  default = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}

/*
When any one of the members approval is required.
*/
variable "account_quorum_policy1" {
  type = any
  description = "The policy document. This is a JSON formatted string."
  default = <<-EOF
         {
          "policy": {
            "quorum": {
              "n": 1,
              "members": [
                {
                  "quorum": {
                    "n": 1,
                    "members": [
                      {
                        "user": "f2fb2f06-3aab-4f76-8e84-638018415db4"
                      }
                    ],
                    "require_2fa": false,
                    "require_password": true
                  }
                },
                {
                  "quorum": {
                    "n": 1,
                    "members": [
                      {
                        "app": "b560a11d-23fb-43e7-ae2d-d6d3c5a1b2ee"
                      }
                    ]
                  }
                }
              ]
            }
          },
          "manage_groups": false,
          "protect_authentication_methods": true,
          "protect_cryptographic_policy": true,
          "protect_logging_config": true
        }
   EOF
}

// When all the members of approval is required.
variable "account_quorum_policy2" {
  type = any
  description = "The policy document. This is a JSON formatted string. First level "
  default = <<-EOF
         {
          "policy": {
            "quorum": {
              "n": 2,
              "members": [
                {
                  "quorum": {
                    "n": 1,
                    "members": [
                      {
                        "user": "f2fb2f06-3aab-4f76-8e84-638018415db4"
                      }
                    ],
                    "require_2fa": false,
                    "require_password": true
                  }
                },
                {
                  "quorum": {
                    "n": 1,
                    "members": [
                      {
                        "app": "b560a11d-23fb-43e7-ae2d-d6d3c5a1b2ee"
                      }
                    ]
                  }
                }
              ]
            }
          },
          "manage_groups": false,
          "protect_authentication_methods": true,
          "protect_cryptographic_policy": true,
          "protect_logging_config": true
        }
   EOF
}


