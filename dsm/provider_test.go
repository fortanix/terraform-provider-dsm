package dsm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestProvider_impl(t *testing.T) {
	var (
		_ *schema.Provider = Provider()
	)
}
