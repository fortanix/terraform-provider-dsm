ChangeLog:
0.5.7:
- using "github.com/hashicorp/go-retryablehttp" as wrapper for http client
- introduce rate limiting from client side
0.5.6:
- fix aws-policy update
- fix aws-tags update
0.5.5:
- Bugfixes
- Support multiple groups for apps
0.5.2:
- GCP SA App bug fix
0.5.0:
- Initial Azure BYOK support
0.4.1:
- Expose all DSM variables
0.4.0:
- Support for App creation and also get credentials
- Add support for deactivation date
0.3.8:
- Support Administrative App API Key as alternative to credentials
0.3.5:
- Moved the /aws_temporary_credential call directly into the NewAPIClient code (seems to be some GPRC issues)
- Fixed issue where alias removal was failing as there was no alias for the CMK
0.3.3: r.statuscode for all APICallBody
0.3.2:
- Clean up at the provider settings level
- Region is now also set at the provider level that is also applied at the AWS credentials support level
0.3.1:
- Moved AWS profile support
0.3.0:
- Removed any AWS KMS group synchronization operations during the BYOK operation
- Created new resource type called "dsm_aws_group" that allows you to dynamically create a group for a specific region
    - attribute "region" should be in short form (such as us-east-1) and is not sanity checked at this very moment - we look to make sure the check is done when v0.3.1 is released
    - access_key / secret_key can be provided at the group creation or can be omitted (this is optional)
- Additional attributes for data block "dsm_aws_group"
    - attribute "scan" is now introduced as a boolean value to scan the AWS KMS group if needed
    - default is "false" for this operation
- Some errors have been updated
    - Not all errors have been meaningful today as the DSM provider relays the error message from the DSM API as-is. This will be continued to fixed as we continue to release the next version
0.2.4: allow specific secret
0.2.3:
Allow insecure SSL communication
0.2.1:
Support secrets allowed to be exported into TF Provider
0.2.0:
GCP EKM support
0.1.9:
Changing the order of scan/check as well as the removal of aws-aliases
0.1.8
Fix KCV issue with RSA keys
Moving to 0.1.8
0.1.5
Added support for Secrets
Renamed all to DSM

How To:
install go
ensure goroot/gopath are set
go get
make
