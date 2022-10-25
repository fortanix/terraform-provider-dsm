# dsm\_sobject\_info

## dsm\_sobject\_info

Returns the DSM security object info from the cluster as a Data Source.

## Usage Reference

```
data "dsm_sobject_info" "sobject" {
    name        = <sobject_name>
}
```

## Argument Reference

The following arguments are supported in the `dsm_sobject_info` data source block:

* **name**: Security object name

## Attribute Reference

The following attributes are stored in the `dsm_sobject_info` data source block:

* **id**: Unique ID of object from Terraform (matches the `kid`)
* **kid**: Security object ID from DSM
* **name**: Security object name from DSM (matches the `name` provided during creation)
* **acct\_id**: Account ID from DSM
* **pub\_key**: Public key from DSM (If applicable)
* **obj\_type**: Security object key type from DSM
* **key\_size**: Security object key size from DSM
* **key\_ops**: Security object key permission from DSM
* **enabled**: true or false
* **creator**: Creator of the security object from DSM
* **description**: Security object description
