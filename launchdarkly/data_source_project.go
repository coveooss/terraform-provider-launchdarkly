package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: resourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKey,
			},
		},
	}
}
