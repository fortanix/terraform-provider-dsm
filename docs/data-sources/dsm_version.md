# dsm\_version

## dsm\_version

Returns the Fortanix DSM version of the cluster as a Data Source.

## Usage Reference

```
data "dsm_version" "version" {}
```

## Argument Reference

None.

## Attribute Reference

The following attributes are stored in the `dsm_version` datasource block:

* **version**: The Fortanix DSM version
* **api\_version**: The Fortanix DSM API version
* **server\_mode**: The Fortanix DSM execution environment
  * **SGX**: The Fortanix DSM running in Intel® SGX environment
