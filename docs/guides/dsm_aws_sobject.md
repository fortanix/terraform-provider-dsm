***Schedule deletion of AWS security object***

**Following attributes should be specified**

*1. pending_window_in_days*

*2. schedule_deletion*

```terraform
// Schedule an DSM AWS security object to delete
/*
Enable schedule_deletion as true and specify pending_window_in_days.
Default value of pending_window_in_days is 7 days.
This can be enabled during both creation and updation. 
*/
resource "dsm_aws_sobject" "dsm_aws_sobject" {
  name = "dsm_aws_sobject"
  group_id = dsm_group.dsm_aws_group.id
  description = "dsm aws sobject"
  key = {
    kid = dsm_sobject.dsm_sobject.id
  }
  enabled = true
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
    aws-policy  = ""
  }
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "APPMANAGEABLE"]
  rotation_policy = {
    interval_days = 10
    effective_at = "20251130T183000Z"
    deactivate_rotated_key = false
  }
  // schedule deletion
  pending_window_in_days = 10
  schedule_deletion = true
}
```
