package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The access token used to authenticate against LaunchDarkly's API",
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"launchdarkly_project":      resourceProject(),
			"launchdarkly_environment":  resourceEnvironment(),
			"launchdarkly_feature_flag": resourceFeatureFlag(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"launchdarkly_project":      dataSourceProject(),
			"launchdarkly_environment":  dataSourceEnvironment(),
			"launchdarkly_feature_flag": dataSourceFeatureFlag(),
		},

		ConfigureFunc: providerConfigure,
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := Client{
		AccessToken: d.Get("access_token").(string),
	}

	return client, nil
}
