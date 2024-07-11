resource "dsm_group" "group1" {
  name = "group1"
  description = "group1"
}

resource "dsm_group" "group2" {
  name = "group2"
  description = "group2"
}

resource "dsm_app" "test_app_create" {
  name = "test_app_create"
  default_group = dsm_group.group1.id
  other_group = [dsm_group.group2.id]
  other_group_permissions = local.other_groups
}


locals {
  other_groups = zipmap(
    [
      dsm_group.group1.id,
      dsm_group.group2.id
    ],
    [
      "SIGN,VERIFY,ENCRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM",
      "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT,TRANSFORM"
    ]
  )
}