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
			"project_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"temporary": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"include_in_snippet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
		},
	}
}

func resourceFeatureFlagCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	description := d.Get("description").(string)
	temporary := d.Get("temporary").(bool)
	includeInSnippet := d.Get("include_in_snippet").(bool)

	payload := map[string]interface{}{
		"name":  name,
		"key":   key,
		"description": description,
		"temporary": temporary,
		"includeInSnippet": includeInSnippet,
	}

	err := client.Post("/flags/"+project, payload, map[int]bool{201: true}, nil)
	if err != nil {
		return err
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)
	d.Set("description", description)
	d.Set("temporary", temporary)
	d.Set("include_in_snippet", includeInSnippet)

	return nil
}

func resourceFeatureFlagRead(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project_key").(string)
	key := d.Get("key").(string)

	client := m.(Client)

	payload := make(map[string]interface{})

	err := client.Get("/flags/"+project+"/"+key, map[int]bool{200: true}, &payload)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", payload["name"])
	d.Set("key", payload["key"])
	d.Set("description", payload["description"])
	d.Set("temporary", payload["temporary"])
	d.Set("include_in_snippet", payload["includeInSnippet"])

	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	temporary := d.Get("temporary").(bool)
	includeInSnippet := d.Get("include_in_snippet").(bool)

	payload := []map[string]interface{}{{
		"op":    "replace",
		"path":  "/name",
		"value": name,
	}, {
		"op":    "replace",
		"path":  "/description",
		"value": description,
	}, {
		"op":    "replace",
		"path":  "/temporary",
		"value": temporary,
	}, {
		"op":    "replace",
			"path":  "/includeInSnippet",
			"value": includeInSnippet,
	}}

	err := client.Patch("/flags/"+project+"/"+d.Id(), payload, map[int]bool{200: true}, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)

	err := client.Delete("/flags/"+project+"/"+d.Id(), map[int]bool{204: true, 404: true})
	if err != nil {
		return err
	}

	return nil
}
