---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dsm_role Data Source - terraform-provider-dsm"
subcategory: ""
description: |-
  Returns the Fortanix DSM role object from the cluster as a Data Source.
---

# dsm_role (Data Source)

Returns the Fortanix DSM role object from the cluster as a Data Source.

## Example Usage

```terraform
data "dsm_role" "sample_role" {
  name = "my_role"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Role name in Fortanix DSM.

### Read-Only

- `id` (String) The ID of this resource.
- `role_id` (String) Role object ID from Fortanix DSM.