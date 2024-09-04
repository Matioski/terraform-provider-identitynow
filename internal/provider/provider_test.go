//go:build !integration

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
)

const (
	providerConfig = `
provider "identitynow" {
  client_id = "clientId"
  client_secret = "clientSecret"
  host     = "http://localhost:3000"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"identitynow": providerserver.NewProtocol6WithError(New("test")()),
	}
	_ = enableLogging()
)

func enableLogging() bool {
	os.Setenv("TF_LOG", "INFO")
	return true
}
