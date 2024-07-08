# dsm\_csrs

## dsm\_csrs

Returns the Fortanix DSM csr sobject from the cluster as a Resource.

## Usage Reference

```
resource "dsm_csr" "sobject_csr" {
    kid   = <sobject_id>
    // Distinguished Name attributes
    cn    = <Common Name for CSR>
    ou    = <Organisational Unit>
    o     = <Organisation>
    l     = <Location>
    c     = <Country>
    st    = <State>
    e     = <Email>
    // Subject Alternative Name(SAN) 
    email = [<Email for CSR>]    
    dnsnames = [<dnsnames>]
    ips   = [<IP addresses>]
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
* _**e**_ : The security object CSR Email
* _**cn**_: The security object CSR Common Name
* _**email**_ : Email value for CSR in Subject Alternative names
* _**dnsnames**_: The security object CSR DNS Names
* _**ips**_: The security object CSR IPs
* _**id**_: The unique ID of object from Terraform

## Note
* Distinguished Name attributes: cn, ou, o, l, c, st and e
* Subject Alternative Name attributes: email, dnsnames and ips