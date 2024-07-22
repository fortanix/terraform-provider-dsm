// Create three normal groups
resource "dsm_group" "group1" {
  name = "group1"
  description = "group1"
}

resource "dsm_group" "group2" {
  name = "group2"
  description = "group2"
}

resource "dsm_group" "group3" {
  name = "group3"
  description = "group3"
}

resource "dsm_gcp_ekm_sa" "sample_ekm_sa" {
  name = "service-[PROJECT-NUMBER]@gcp-sa-ekms.iam.gserviceaccount.com"
  default_group = dsm_group.group1.id
  other_group =     [
      dsm_group.group2.id,
      dsm_group.group3.id
    ]

  other_group_permissions = zipmap(
    [
      dsm_group.group2.id,
      dsm_group.group3.id
    ],
    [
      "SIGN,VERIFY,ENCRYPT,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT",
      "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT"
    ]
  )
}