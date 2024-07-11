# dsm\_aws\_group

## dsm\_aws\_group

Returns the Fortanix DSM AWS KMS mapped group object from the cluster as a Data Source for AWS KMS.

## Usage Reference

```
data "dsm_aws_group" "aws_group" {
    name = <group_name>
    scan = <true/false>
}
```

## Argument Reference

The following arguments are supported in the `dsm_aws_group` data source block:

* **name**: The AWS KMS group object name in Fortanix DSM
* _**scan (optional)**_: Syncs keys from AWS KMS to the AWS KMS group in DSM. Value is either `True`/`False`

## Attribute Reference

The following attributes are stored in the `dsm_aws_group` data source block:

* **id**: The unique ID of object from Terraform (matches the `group_id`) 
* **group\_id**: The AWS KMS group object ID from Fortanix DSM
* **name**: The AWS KMS group object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The Account ID from Fortanix DSM
* **region**: The AWS region mapped to the group from which keys are imported
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The AWS KMS group object description

In addition, the following attributes will be used to communicate with the corresponding AWS KMS instance:

* **access\_key**: The Access Key ID used to communicate with AWS KMS
