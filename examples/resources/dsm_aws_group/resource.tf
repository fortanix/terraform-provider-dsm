// Creation of an dsm aws group
// Default aws_region is us-east-1. It can be modified by adding it in provider.
provider "dsm" {
  aws_region = "us-east-2"
}

resource "dsm_aws_group" "dsm_aws_group_terraform" {
  name = "dsm_aws_group_terraform"
  description = "AWS group"
  access_key = "XXXXXXXXXXXXXXXXXXXX"
  secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}