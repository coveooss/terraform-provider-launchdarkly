package main

import (
	"github.com/coveo/terraform-provider-launchdarkly/launchdarkly"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return launchdarkly.Provider()
		},
	})
}
