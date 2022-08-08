# dsm\_sobject

## dsm\_sobject

Returns the DSM security object from the cluster as a Data Source.

## Usage Reference

```
data "dsm_sobject" "sobject" {
    name        = <sobject_name>
}
```

## Argument Reference

The following arguments are supported in the `dsm_sobject` data source block:

* **name**: Security object name

## Attribute Reference

The following attributes are stored in the `dsm_sobject` data source block:

* **id**: Unique ID of object from Terraform (matches the `kid`)
* **kid**: Security object ID from DSM
* **name**: Security object name from DSM (matches the `name` provided during creation)
* **export**: true or false
* **acct\_id**: Account ID from DSM
* **obj\_type**: Security object key type from DSM
* **key\_size**: Security object key size from DSM
* **key\_ops**: Security object key permission from DSM
* **enabled**: true or false
* **value**: Value of key material (only if export is allowed)
* **creator**: Creator of the security object from DSM
  * **user**: If the security object was created by a user, the computed value will be the matching user id
  * **app**: If the security object was created by a app, the computed value will be the matching app id
* **description**: Security object description
