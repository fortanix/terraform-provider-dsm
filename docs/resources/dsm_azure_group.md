# dsm\_azure\_group

## dsm\_azure\_group

Returns the Fortanix DSM group mapped to Azure KV  from the cluster as a resource.

## Usage Reference

```
resource "dsm_azure_group" "azure_group" {
    name            = <group_name>
    description     = <description of the group>
    url             = <Azure Key Vault URL>
    client_id       = <Azure App ID>
    subscription_id = <Azure Subscription ID>
    tenant_id       = <Azure Tenant ID>
    secret_key      = <Azure App secret>
    key_vault_type  = <Standard/Premium>
}
```

## Argument Reference

The following arguments are supported in the `dsm_azure_group` resource block:

* **name**: The Azure KV group object name in Fortanix DSM
* _**description (optional)**_: Description of the Azure KV Fortanix DSM group
* **URL**: The URL of the object in an Azure KV that uniquely identifies the object
* **client_id**: The Azure registered application id (username)
* **subscription\_id**: The ID of the Azure AD subscription
* **tenant\_id**: The tenant/directory id of the Azure subscription
* _**key\_vault\_type (optional)**_: The type of key vaults, Standard/Premium
* **secret\_key**: A secret string that a registered application in Azure uses to prove its identity (application password) 

## Attribute Reference

The following attributes are stored in the `dsm_azure_group` resource block:

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
