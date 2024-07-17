resource "dsm_gcp_ekm_sa" "sample_ekm_sa" {
  name = "service-[PROJECT-NUMBER]@gcp-sa-ekms.iam.gserviceaccount.com"
  default_group = "035f84b5-75b6-4f37-8968-19d588695bcc"
  other_group =     [
      "d23ea001-80ef-41ad-a3b2-25489b555a58",
      "499760d3-04bd-4c7e-98b1-51a9e8871ee1"
    ]

  other_group_permissions = zipmap(
    [
      "d23ea001-80ef-41ad-a3b2-25489b555a58",
      "499760d3-04bd-4c7e-98b1-51a9e8871ee1"
    ],
    [
      "SIGN,VERIFY,ENCRYPT,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT",
      "SIGN,VERIFY,DECRYPT,WRAPKEY,UNWRAPKEY,DERIVEKEY,MACGENERATE,MACVERIFY,EXPORT,MANAGE,AGREEKEY,AUDIT"
    ]
  )
}