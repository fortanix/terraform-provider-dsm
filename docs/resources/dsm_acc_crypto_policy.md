# dsm\_acc\_crypto\_policy

## dsm\_acc\_crypto\_policy

Returns the Fortanix DSM account cryptographic policy object from the cluster as a Resource.

## Usage Reference

```
resource "dsm_acc_crypto_policy" "name" {
    acct_id = <account_id>
    cryptographic_policy = <account_cryptographic_policy_definition>
}
```

## Argument Reference

The following arguments are supported and required in the `dsm_acc_crypto_policy` resource block:

* **acct\_id**: The Fortanix DSM account object id.
* **cryptographic\_policy**: The Fortanix DSM account object cryptographic policy definition as a JSON string

## Attribute Reference

The following attributes are stored in the `dsm_acc_crypto_policy` resource block:

* **id**: Unique ID of object from Terraform (matches the `acct_id` from resource block)
* **acct\_id**: Account ID from Fortanix DSM
* **approval\_policy**: The Fortanix DSM account object quorum approval policy definition as a JSON string
* **cryptographic\_policy**: The Fortanix DSM account object cryptographic policy definition as a JSON string