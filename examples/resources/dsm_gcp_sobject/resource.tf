resource "dsm_gcp_sobject" "sample_gcp_sobject" {
    name = "test-gcp-sobject"
    group_id = "311915f2-7cdb-4ea9-ac15-83818f04dc39"
    key = {
      kid = "f6be4755-7912-4546-9e94-27851f2ddcd7"
    }
    custom_metadata = {
      gcp-key-id = "name-of-the-key-in-gcp"
    }
}