# dsm\_azure\_group

## dsm\_azure\_group

Returns the Fortanix DSM Azure KV mapped group object from the cluster as a Data Source for Azure Key Vault.

## Usage Reference

```
data "dsm_azure_group" "azure_group" {
    name = <group_name>
    scan = <true/false>
}
```

## Argument Reference

The following arguments are supported in the `dsm_azure_group` data source block:

* **name**: The Azure KV group object name in Fortanix DSM
* _**scan (optional)**_: Syncs keys from Azure KV to the Azure KV group in Fortanix DSM. Value is either `True`/`False`

## Attribute Reference

The following attributes are stored in the `dsm_azure_group` data source block:

* **id**: The unique ID of object from Terraform (matches the `group_id`) 
* **group\_id**: The Azure KV group object ID from Fortanix DSM
* **name**: The Azure KV group object name from Fortanix DSM (matches the name provided during creation)
* **acct\_id**: The Account ID from Fortanix DSM
* **creator**: The creator of the group object from Fortanix DSM
  * **user**: If the group was created by a user, the computed value will be the matching user id
  * **app**: If the group was created by a app, the computed value will be the matching app id
* **description**: The Azure KV group object description
* **url**: The URL of the object in an Azure KV that uniquely identifies the object
* **client\_id**: The Azure registered application id (username)
*	**subscription\_id**: The ID of the Azure AD subscription
*	**tenant\_id**: The tenant/directory id of the Azure subscription
*	**key\_vault\_type**: The type of key vaults, Standard/Premium 
