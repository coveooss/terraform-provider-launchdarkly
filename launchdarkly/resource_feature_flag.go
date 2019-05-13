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
		Importer: &schema.ResourceImporter{
			State: resourceFeatureFlagImport,
		},

		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKey,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateFeatureFlagKey,
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

func resourceFeatureFlagImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectKey, environment, err := parseCompositeID(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(environment)
	d.Set("project_key", projectKey)
	d.Set("key", environment)

	resourceFeatureFlagRead(d, meta)

	return []*schema.ResourceData{d}, nil
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

	payload := JsonFeatureFlag{
		Name:             name,
		Key:              key,
		Description:      description,
		Temporary:        temporary,
		IncludeInSnippet: includeInSnippet,
		Tags:             transformTagsFromTerraformFormat(tags),
		CustomProperties: transformedCustomProperties,
	}

	var response JsonFeatureFlag
	err = client.Post(getFlagCreateUrl(project), payload, []int{201}, &response)
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

	var response JsonFeatureFlag
	err := client.GetInto(getFlagUrl(project, key), []int{200}, &response)
	if err != nil {
		d.SetId("")
		return nil
	}

	transformedCustomProperties := transformCustomPropertiesFromLaunchDarklyFormat(response.CustomProperties)

	d.SetId(key)
	d.Set("name", response.Name)
	d.Set("key", response.Key)
	d.Set("description", response.Description)
	d.Set("temporary", response.Temporary)
	d.Set("include_in_snippet", response.IncludeInSnippet)
	d.Set("tags", response.Tags)
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

	_, err = client.Patch(getFlagUrl(project, d.Id()), payload, []int{200})
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)

	err := client.Delete(getFlagUrl(project, d.Id()), []int{204, 404})
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

func transformCustomPropertiesFromTerraformFormat(properties []interface{}) (map[string]JsonCustomProperty, error) {
	transformed := make(map[string]JsonCustomProperty)

	for _, raw := range properties {
		value := raw.(map[string]interface{})
		key := value["key"].(string)
		name := value["name"].(string)

		values := []string{}
		for _, v := range value["value"].([]interface{}) {
			values = append(values, v.(string))
		}

		transformed[key] = JsonCustomProperty{
			Name:  name,
			Value: values,
		}
	}

	return transformed, nil
}

func transformCustomPropertiesFromLaunchDarklyFormat(properties map[string]JsonCustomProperty) interface{} {
	transformed := make([]map[string]interface{}, 0)

	for key, body := range properties {
		sub := make(map[string]interface{})
		sub["key"] = key
		sub["name"] = body.Name
		sub["value"] = body.Value

		transformed = append(transformed, sub)
	}

	return transformed
}
