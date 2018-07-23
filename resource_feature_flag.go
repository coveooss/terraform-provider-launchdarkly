package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"errors"
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
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
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
	tags := d.Get("tags").([]interface{})
	customProperties := d.Get("custom_properties").(map[string]interface{})

	transformedCustomProperties, err := transformCustomPropertiesFromTerraformFormat(customProperties)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":  name,
		"key":   key,
		"description": description,
		"temporary": temporary,
		"includeInSnippet": includeInSnippet,
		"tags": transformTags(tags),
		"customProperties": transformedCustomProperties,
	}

	err = client.Post("/flags/"+project, payload, map[int]bool{201: true}, nil)
	if err != nil {
		return err
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)
	d.Set("description", description)
	d.Set("temporary", temporary)
	d.Set("include_in_snippet", includeInSnippet)
	d.Set("tags", tags)
	d.Set("custom_properties", customProperties)

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

	transformedTags := payload["tags"].([]interface{})
	transformedCustomProperties := transformCustomPropertiesFromLaunchDarklyFormat(payload["customProperties"])

	d.Set("name", payload["name"])
	d.Set("key", payload["key"])
	d.Set("description", payload["description"])
	d.Set("temporary", payload["temporary"])
	d.Set("include_in_snippet", payload["includeInSnippet"])
	d.Set("tags", transformedTags)
	d.Set("custom_properties", transformedCustomProperties)

	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	temporary := d.Get("temporary").(bool)
	includeInSnippet := d.Get("include_in_snippet").(bool)
	tags := d.Get("tags").([]interface{})
	customProperties := d.Get("custom_properties").(map[string]interface{})

	transformedCustomProperties, err := transformCustomPropertiesFromTerraformFormat(customProperties)
	if err != nil {
		return err
	}

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
	}, {
		"op":    "replace",
		"path":  "/tags",
		"value": transformTags(tags),
	}, {
		"op":    "replace",
		"path":  "/customProperties",
		"value": transformedCustomProperties,
	}}

	err = client.Patch("/flags/"+project+"/"+d.Id(), payload, map[int]bool{200: true}, nil)
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

func transformTags(tags []interface{}) []string {
	transformed := make([]string, len(tags))

	for index, value := range tags {
		transformed[index] = value.(string)
	}

	return transformed
}

func transformCustomPropertiesFromTerraformFormat(properties map[string]interface{}) (interface{}, error) {
	transformed := make(map[string]interface{})

	for key, value := range properties {
		sub := make(map[string]interface{})
		sub["name"] = key

		switch value.(type) {
		case string:
			sub["value"] = []string { value.(string) }
		default:
			return transformed, errors.New("Custom property " + key + " must have value of type string")
		}

		transformed[key] = sub
	}

	return transformed, nil
}

func transformCustomPropertiesFromLaunchDarklyFormat(properties interface{}) interface{} {
	transformed := make(map[string]interface{})

	for key, body := range properties.(map[string]interface{}) {
		values := body.(map[string]interface{})["value"].([]interface{})

		if len(values) != 1 {
			println("Skipping custom property " + key + " because it has multiple values")
			continue
		}

		transformed[key] = values[0].(string)
	}

	return transformed
}
