# Rotation of dsm_aws_sobject

## Rotate with DSM Option

```terraform
resource "dsm_group" "normal_group" {
  name = "normal_group"
}

// Create AWS group
resource "dsm_group" "aws_group" {
  name = "aws_group"
  description = "AWS group"
  hmg = jsonencode(
    {
      url = "kms.us-east-1.amazonaws.com"
      tls = {
        mode = "required"
        validate_hostname: false,
        ca = {
          ca_set = "global_roots"
        }
      }
      kind = "AWSKMS"
      access_key = "XXXXXXXXXXXXXXXXXXXX"
      secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
      region = "us-east-1"
      service = "kms"
    })
}

// Create a dsm_sobject of type AES key inside DSM
resource "dsm_sobject" "aes_sobject" {
  name            = "aes_sobject"
  obj_type        = "AES"
  group_id        = dsm_group.normal_group.id
  key_size        = 256
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
}

// AWS sobject creation(Copies the key from DSM)
resource "dsm_aws_sobject" "aws_sobject" {
  name = "aws_sobject"
  group_id = dsm_group.aws_group.id
  description = "AWS sobject"
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  key = {
    kid = dsm_sobject.aes_sobject.id
  }
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
  }
}

// 1st Rotation of dsm_sobject
resource "dsm_sobject" "aes_sobject_rotate1" {
  name            = "aes_sobject"
  obj_type        = "AES"
  group_id        = dsm_group.normal_group.id
  key_size        = 256
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  rotate          = "DSM"
  rotate_from     = dsm_sobject.aes_sobject.name // Name of the old dsm_sobject
  
}

// 1st Rotation of dsm_aws_sobject
resource "dsm_aws_sobject" "aws_sobject_rotate1" {
  name = "aws_sobject"
  group_id = dsm_group.aws_group.id
  description = "AWS sobject"
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  key = {
    kid = dsm_sobject.aes_sobject_rotate1.id // 1st Rotated dsm_sbject
  }
  rotate = "DSM"
  rotate_from = dsm_aws_sobject.aws_sobject.name // Name of the old dsm_aws_sobject
}

// 2nd Rotation of dsm_sobject
resource "dsm_sobject" "aes_sobject_rotate2" {
  name            = "aes_sobject"
  obj_type        = "AES"
  group_id        = dsm_group.normal_group.id
  key_size        = 256
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  rotate          = "DSM"
  rotate_from     = dsm_sobject.aes_sobject1.name // Name of the old dsm_sobject

}

// 2nd Rotation of dsm_aws_sobject
resource "dsm_aws_sobject" "aws_sobject_rotate2" {
  name = "aws_sobject"
  group_id = dsm_group.aws_group.id
  description = "AWS sobject"
  key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
  key = {
    kid = dsm_sobject.aes_sobject_rotate2.id // 2nd Rotated dsm_sbject
  }
  rotate = "DSM"
  rotate_from = dsm_aws_sobject.aws_sobject_rotate1.name // Name of the old dsm_aws_sobject
}
```

## Rotate with AWS Option

```terraform
// 1st Rotation of dsm_sobject
resource "dsm_sobject" "aes_sobject_rotate1" {
name            = "aes_sobject"
obj_type        = "AES"
group_id        = dsm_group.normal_group.id
key_size        = 256
key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
rotate          = "DSM"
rotate_from     = dsm_sobject.aes_sobject.name // Name of the old dsm_sobject

}

// 1st Rotation of dsm_aws_sobject
resource "dsm_aws_sobject" "aws_sobject_rotate1" {
name = "aws_sobject"
group_id = dsm_group.aws_group.id
description = "AWS sobject"
key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
key = {
kid = dsm_sobject.aes_sobject_rotate1.id // 1st Rotated dsm_sbject
}
rotate = "AWS"
rotate_from = dsm_aws_sobject.aws_sobject.name // Name of the old dsm_aws_sobject
}

// 2nd Rotation of dsm_sobject
resource "dsm_sobject" "aes_sobject_rotate2" {
name            = "aes_sobject"
obj_type        = "AES"
group_id        = dsm_group.normal_group.id
key_size        = 256
key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
rotate          = "DSM"
rotate_from     = dsm_sobject.aes_sobject1.name // Name of the old dsm_sobject

}

// 2nd Rotation of dsm_aws_sobject
resource "dsm_aws_sobject" "aws_sobject_rotate2" {
name            = "aws_sobject"
group_id        = dsm_group.aws_group.id
description     = "AWS sobject"
key_ops         = ["EXPORT", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "APPMANAGEABLE"]
key             = {
                  kid = dsm_sobject.aes_sobject_rotate2.id // 2nd Rotated dsm_sbject
                  }
rotate          = "AWS"
rotate_from     = dsm_aws_sobject.aws_sobject_rotate1.name // Name of the old dsm_aws_sobject
}

```


# Schedule deletion and Delete Key Material of an AWS security object

## Following attributes should be specified for schedule deletion

*1. schedule_deletion*

```terraform
// Schedule an DSM AWS security object to delete
/*
Enable schedule_deletion as an Integer value.
The minimum value of a schedule_deletion is 7days.
This can be enabled only during update. 
*/
resource "dsm_aws_sobject" "dsm_aws_sobject" {
  name = "dsm_aws_sobject"
  group_id = dsm_group.dsm_aws_group.id
  description = "dsm aws sobject"
  key = {
    kid = dsm_sobject.dsm_sobject.id
  }
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
  }
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "APPMANAGEABLE"]
  rotation_policy = {
    interval_days = 10
    effective_at = "20251130T183000Z"
    deactivate_rotated_key = false
  }
  // schedule deletion
  schedule_deletion = 7
}
```

## Following attributes should be specified for Delete key material

*1. delete_key_material*

```terraform
/*
Enable delete_key_material as true.
This can be enabled only during update. 
*/
resource "dsm_aws_sobject" "dsm_aws_sobject" {
  name = "dsm_aws_sobject"
  group_id = dsm_group.dsm_aws_group.id
  description = "dsm aws sobject"
  key = {
    kid = dsm_sobject.dsm_sobject.id
  }
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
  }
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "APPMANAGEABLE"]
  rotation_policy = {
    interval_days = 10
    effective_at = "20251130T183000Z"
    deactivate_rotated_key = false
  }
  // delete key material
  delete_key_material = true
}
```

## Both schedule_deletion and delete_key_material can be enabled in a single terraform request

```terraform
/*
Enable delete_key_material as true.
This can be enabled only during update. 
*/
resource "dsm_aws_sobject" "dsm_aws_sobject" {
  name = "dsm_aws_sobject"
  group_id = dsm_group.dsm_aws_group.id
  description = "dsm aws sobject"
  key = {
    kid = dsm_sobject.dsm_sobject.id
  }
  custom_metadata = {
    aws-aliases = "dsm_aws_sobject"
  }
  key_ops         = ["ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "APPMANAGEABLE"]
  rotation_policy = {
    interval_days = 10
    effective_at = "20251130T183000Z"
    deactivate_rotated_key = false
  }
  // schedule_deletion and delete key material
  schedule_deletion = 7
  delete_key_material = true
}
```