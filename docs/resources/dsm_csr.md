# dsm\_csrs

## dsm\_csrs

Returns the Fortanix DSM csr sobject from the cluster as a Resource.

## Usage Reference

```
resource "dsm_csr" "sobject" {
    kid   = <sobject_id>
    cn    = <Common Name for CSR>
    email = <Email for CSR>    
}
```

## Argument Reference

The following arguments are supported in the `dsm_csrs` resource block:

* _**kid**_ : The security object kid
* _**value**_ : The security object value of Generated CSR
* _**ou**_ : The security object CSR Organisational Unit
* _**o**_ : The security object CSR Organisation
* _**l**_ : The security object CSR Location
* _**c**_ : The security object CSR Country
* _**st**_ :  The security object CSR State
* _**email**_ : Email value for CSR
* _**cn**_: The security object CSR Common Name
* _**dnsnames**_: The security object CSR DNS Names
* _**ips**_: The security object CSR IPs
* _**id**_: The unique ID of object from Terraform