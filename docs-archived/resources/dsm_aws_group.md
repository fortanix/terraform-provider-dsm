# dsm\_aws\_group

## dsm\_aws\_group

Returns Fortanix DSM group mapped to AWS KMS from the cluster as a resource.

## Usage Reference

```
resource "dsm_aws_group" "aws_group" {
    name        = <Custom Group Name>-AWS-<Region>    
    access_key  = <access_key>
    secret_key  = <secret_key>
}
```

## Argument Reference

The following arguments are supported in the `dsm_aws_group` resource block:

* **name**: The name follows the nomenclature of `<Custom Group Name>-AWS-<Region>`
* _**description (optional)**_: The description of the AWS KMS group
* _**access\_key (optional)**_: Th Access Key ID to set for AWS KMS group for programmatic (API) access to AWS Services
* _**secret\_key (optional)**_: The Secret Access Key to set for AWS KMS group for programmatic (API) access to AWS Services

## Attribute Reference

The following attributes are stored in the `dsm_aws_group` resource block:

* **id**: The unique ID of object from Terraform (matches the `group_id`)
* **name**: The Fortanix DSM AWS KMS mapped group Name (matches the name provided during creation) 
* **group\_id**: The unique ID for AWS KMS Mapped group from Fortanix DSM
* **acct\_id**: The Account ID from Fortanix DSM
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **region**: The AWS region mapped to the group from which keys are imported
* **description**: The AWS KMS group object description
* **access\_key**: The Access Key ID used to communicate with AWS KMS
