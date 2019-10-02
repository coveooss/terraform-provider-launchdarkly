package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFeatureFlag() *schema.Resource {
	return &schema.Resource{
		Read: resourceFeatureFlagRead,

		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKey,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"temporary": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"include_in_snippet": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}
