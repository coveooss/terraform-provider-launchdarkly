package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFeatureFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureFlagCreate,
		Read:   resourceFeatureFlagRead,
		Update: resourceFeatureFlagUpdate,
		Delete: resourceFeatureFlagDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFeatureFlagCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFeatureFlagRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
