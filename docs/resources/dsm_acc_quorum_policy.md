# dsm\_acc\_quorum\_policy

## dsm\_acc\_quorum\_policy

Returns the Fortanix DSM account quorum policy object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_acc_quorum_policy" "policyname" {
    acct_id = <account_id>
    approval_policy = <account_policy_description>
}
```

## Argument Reference

The following arguments are supported in the `dsm_group` resource block:

* **acct_id**: The Fortanix DSM account object id.
* _**approval_policy (optional)**_: The Fortanix DSM account object quorum approval policy definition as a JSON string

## Attribute Reference

The following attributes are stored in the `dsm_group` resource block:

* **id**: Unique ID of object from Terraform (matches the `group_id` from resource block)
* **acct\_id**: Account ID from Fortanix DSM
* **approval_policy**: The Fortanix DSM account object quorum approval policy definition as a JSON string

## Note

Since modifying or deleting an already existing quorum policy will require approval, these operations are not meaningful in an automation. Therefore, only create operation is supported.
Modifying or deleting quorum policies should be performed from the UI.
