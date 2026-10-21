package dsm

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"regexp"
	"testing"
)

const (
	ProviderName = "dsm"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
	//testAccProviderEndpoint string = "https://sdkms.fortanix.com"
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"dsm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Attributes whose names look like a secret must be marked Sensitive so Terraform redacts them in plan/apply output.
// Note (Shashidhar Naraparaju): If a matching attribute is not actually a secret, an allowlist map can be added as and when required.
var secretAttributeName = regexp.MustCompile(`credential|secret|password|token|private_key|api_key|access_key`)

func TestSensitiveAttributes(t *testing.T) {
	p := Provider()
	checkSensitive(t, "provider", p.Schema)
	for name, r := range p.ResourcesMap {
		checkSensitive(t, "resource."+name, r.Schema)
	}
	for name, r := range p.DataSourcesMap {
		checkSensitive(t, "data."+name, r.Schema)
	}
}

func checkSensitive(t *testing.T, path string, attrs map[string]*schema.Schema) {
	for name, s := range attrs {
		if s.Sensitive || s.Type == schema.TypeBool {
			continue
		}
		if r, ok := s.Elem.(*schema.Resource); ok {
			checkSensitive(t, path+"."+name, r.Schema)
		} else if secretAttributeName.MatchString(name) {
			t.Errorf("%s.%s looks like a secret but is not marked Sensitive", path, name)
		}
	}
}

func testAccPreCheck(t *testing.T) {
	for _, env := range []string{"DSM_ENDPOINT", "DSM_ACCT_ID", "DSM_USERNAME", "DSM_PASSWORD"} {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

func testAccPreCheckAws(t *testing.T) {
	for _, env := range []string{"AWS_ACCESS_KEY", "AWS_SECRET_KEY"} {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("%s environment variable must be set for AWS BYOK tests", env)
		}
	}
}

func testAccPreCheckGcp(t *testing.T) {
	for _, env := range []string{"GOOGLE_SERVICE_ACCOUNT"} {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("%s environment variable must be set for GCP EKM tests", env)
		}
	}
}

func testAccPreCheckAzure(t *testing.T) {
	for _, env := range []string{"AZURE_TENANT_ID", "AZURE_SECRET_KEY", "AZURE_SUBSCRIPTION_ID", "AZURE_CLIENT_ID", "AZURE_URL"} {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("%s environment variable must be set for Azure Key Vault tests", env)
		}
	}
}

func TestProvider_impl(t *testing.T) {
	var (
		_ *schema.Provider = Provider()
	)
}
