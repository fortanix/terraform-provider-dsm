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
* _**ou**_ : The security object cert Organisational Unit
* _**o**_ : The security object cert Organisation
* _**l**_ : The security object cert Location
* _**c**_ : The security object cert Country
* _**st**_ :  The security object cert State
* _**email**_ : Email value for cert
* _**cn**_: The security object cert Common Name
* _**dnsnames**_: The security object cert DNS Names
* _**ips**_: The security object cert IPs
* _**id**_: The unique ID of object from Terraform