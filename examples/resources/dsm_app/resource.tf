# Create three normal groups
resource "dsm_group" "group1" {
  name        = "group1"
  description = "group1"
}

resource "dsm_group" "group2" {
  name        = "group2"
  description = "group2"
}

resource "dsm_group" "group3" {
  name        = "group3"
  description = "group3"
}

# Create an app
resource "dsm_app" "app" {
  name                    = "app"
  default_group           = dsm_group.group1.id
  other_group             = [dsm_group.group2.id, dsm_group.group3.id]
  other_group_permissions = local.other_groups
}


######################################################################################################################
# To modify the default group permissions, other_group_permissions can be used in a zipmap.
# In the above app the following permissions are assigned for each group:
# group1: SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT
# group2: SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT
# group3: SIGN,VERIFY,ENCRYPT,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT

# For group3, default permissions are assigned as it was not specified in the other_group_permissions.
# group should be specified only if default permissions need to be changed.
######################################################################################################################

locals {
  other_groups = zipmap(
    [
      dsm_group.group1.id,
      dsm_group.group2.id
    ],
    [
      "SIGN,VERIFY,ENCRYPT,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT",
      "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT"
    ]
  )
}

# An example on how to modify the existing permissions of a group in app
resource "dsm_app" "app" {
  name                    = "app"
  default_group           = dsm_group.group1.id
  other_group             = [dsm_group.group2.id, dsm_group.group3.id]
  other_group_permissions = local.other_groups
  # mod_group_permissions should be given while updating an app
  mod_group_permissions = local.mod_groups
}

# group1 and group2 permissions modification
locals {
  mod_groups = zipmap(
    [
      dsm_group.group1.id,
      dsm_group.group2.id
    ],
    [
      "SIGN,VERIFY,ENCRYPT,DECRYPT,WRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT",
      "SIGN,VERIFY,DECRYPT,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT"
    ]
  )
}

