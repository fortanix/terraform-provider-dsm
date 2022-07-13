# dsm\_csrs

## dsm\_csrs

Returns the Fortanix DSM cert sobject object details from the cluster as a Data Source.

## Usage Reference

```
data "dsm_csr" "group" {
    kid = <>
    cn = <>
}
```

## Argument Reference

The following arguments are supported in the `dsm_cert` resource block:

* _**kid**_ : Sobject key id value
* _**cn**_ : Certificate Common Name

## Attribute Reference

The following attributes are stored in the `dsm_cert` data source block:

* _**kid**_ : The security object kid
* _**value**_ : The security object value of Generated CSR
* _**ou**_ : The security object cert Organisational Unit
* _**o**_ : The security object cert Organisation
* _**l**_ : The security object cert Location
* _**c**_ : The security object cert Country
* _**cn**_: The security object cert Common Name
* _**id**_: The unique ID of object from Terraform
