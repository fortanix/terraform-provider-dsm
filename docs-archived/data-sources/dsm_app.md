# dsm\_app

## dsm\_app

Returns the Fortanix DSM app object from the cluster as a Data Source.

## Usage Reference

```
data "dsm_app" "app" {
    app_id = <app_id>
}
```

## Argument Reference

The following arguments are supported in the `dsm_app` data block:

* **app_id**: App id value

## Attribute Reference

The following attributes are stored in the `dsm_app` data source block:

* **id**: The unique ID of object from Terraform
* **credential**: The Fortanix DSM App API key
* _**new\_credential**_: Set to false by default


