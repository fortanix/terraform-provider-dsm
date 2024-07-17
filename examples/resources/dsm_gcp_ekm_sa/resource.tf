resource "dsm_gcp_ekm_sa" "sample_ekm_sa" {
  name = "service-[PROJECT-NUMBER]@gcp-sa-ekms.iam.gserviceaccount.com"
  default_group = "035f84b5-75b6-4f37-8968-19d588695bcc"
}