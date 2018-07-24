package launchdarkly

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
			"project_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"temporary": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"include_in_snippet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
	customProperties := d.Get("custom_properties").([]interface{})

	transformedCustomProperties, err := transformCustomPropertiesFromTerraformFormat(customProperties)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":             name,
		"key":              key,
		"description":      description,
		"temporary":        temporary,
		"includeInSnippet": includeInSnippet,
		"tags":             transformTagsFromTerraformFormat(tags),
		"customProperties": transformedCustomProperties,
	}

	_, err = client.Post(getFlagCreateUrl(project), payload, map[int]bool{201: true})
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

	raw, err := client.Get(getFlagUrl(project, key), map[int]bool{200: true})
	if err != nil {
		d.SetId("")
		return nil
	}

	payload := raw.(map[string]interface{})
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
	customProperties := d.Get("custom_properties").([]interface{})

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
		"value": transformTagsFromTerraformFormat(tags),
	}, {
		"op":    "replace",
		"path":  "/customProperties",
		"value": transformedCustomProperties,
	}}

	_, err = client.Patch(getFlagUrl(project, d.Id()), payload, map[int]bool{200: true})
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)

	err := client.Delete(getFlagUrl(project, d.Id()), map[int]bool{204: true, 404: true})
	if err != nil {
		return err
	}

	return nil
}

func transformTagsFromTerraformFormat(tags []interface{}) []string {
	transformed := make([]string, len(tags))

	for index, value := range tags {
		transformed[index] = value.(string)
	}

	return transformed
}

func transformCustomPropertiesFromTerraformFormat(properties []interface{}) (interface{}, error) {
	transformed := make(map[string]interface{})

	for _, value := range properties {
		key := value.(map[string]interface{})["key"].(string)
		name := value.(map[string]interface{})["name"].(string)
		values := value.(map[string]interface{})["value"].([]interface{})

		sub := make(map[string]interface{})
		sub["name"] = name
		sub["value"] = values

		transformed[key] = sub
	}

	return transformed, nil
}

func transformCustomPropertiesFromLaunchDarklyFormat(properties interface{}) interface{} {
	transformed := make([]map[string]interface{}, 0)

	for key, body := range properties.(map[string]interface{}) {
		name := body.(map[string]interface{})["name"].(string)
		values := body.(map[string]interface{})["value"].([]interface{})

		sub := make(map[string]interface{})
		sub["key"] = key
		sub["name"] = name
		sub["value"] = values

		transformed = append(transformed, sub)
	}

	return transformed
}
