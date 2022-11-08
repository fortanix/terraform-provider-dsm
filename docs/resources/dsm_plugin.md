# dsm\_plugin

## dsm\_plugin

Returns the Fortanix DSM plugin object from the cluster as a Resource.

## Usage Reference

```
data "local_file" "plugin_code" {
    filename ="<file_path>"
}

resource "dsm_plugin" "plugin" {
	name          = "<plugin_name>"
	description   = "<plugin_description>"
	default_group = "<default_group_id>"
	groups        = ["<default_group_id>", "<group_id_1>", ...]
	language      = "<plugin_code_language>"
    code          = data.local_file.plugin_code.content
    plugin_type   = "<plugin_type>"
    enabled       = "<true/false>"
}
```

## Argument Reference

The following arguments are supported in the `dsm_group` resource block:

* **name**: The Fortanix DSM plugin object name
* _**description (optional)**_: The Fortanix DSM plugin object description
* **default_group**: The Fortanix DSM group object id to be mapped to the plugin by default
* **groups**: The Fortanix DSM group object ids to be mapped to the plugin
* _**language (optional)**_: Programming language for plugin code (Default value is  "LUA")
* **code**: Plugin code that will be executed in DSM. Code should be in above given programming language.
* _**plugin_type (Optional)**_: Type of the plugin
* * _**enabled (Optional)**: Whether this plugin is enabled

## Note
Argument code can be given as a string or tag to a file. Please refer the below example to tag a file.

    ```
    Example:
    A local_file variable should be created and filepath should be given as a filename inside local_file.
    Then tag this local variable to the code. 

    data "local_file" "plugin_code" {
        filename ="<file_path>"
    }

    code = data.local_file.plugin_code.content
    ```


## Attribute Reference

The following attributes are stored in the `dsm_plugin` resource block:

* **id**: Unique ID of object from Terraform (matches the `plugin_id` from resource block)
* **plugin\_id**: Plugin object ID from Fortanix DSM
* **name**: Plugin object name from Fortanix DSM (matches the `name` provided during creation)
* **acct\_id**: Account ID from Fortanix DSM
* **creator**: Creator of the group object from Fortanix DSM
    * **user**: If the plugin was created by a user, the computed value will be the matching user id
* **description**: The Fortanix DSM plugin object description
* **default_group**: The default group id mapped to the Fortanix DSM plugin
* **groups**: Group ids to be mapped to the Fortanix DSM plugin
* **language**: Programming language for plugin code
* **code**: Plugin code that will be executed in DSM
* **plugin_type**: Type of the plugin
* **enabled**: Whether this plugin is enabled
* * **approval\_request\_id**: If a plugin creation requires approval, then request id will be stored here.