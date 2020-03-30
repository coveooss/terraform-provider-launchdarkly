package launchdarkly

import (
	"errors"
	"fmt"
	"strconv"
	"sort"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

const VARIATION_NAME_KEY = "name"
const VARIATION_DESCRIPTION_KEY = "description"
const VARIATION_VALUE_KEY = "value"
const VARIATIONS_STRING_KIND = "string"
const VARIATIONS_NUMBER_KIND = "number"
const VARIATIONS_BOOLEAN_KIND = "boolean"
const DEFAULT_VARIATIONS_KIND = VARIATIONS_BOOLEAN_KIND
const NUMBER_OF_RETRY = 3

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
			"variations_kind": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DEFAULT_VARIATIONS_KIND,
				ValidateFunc: validateFeatureFlagVariationsType,
				ForceNew:     true,
			},
			"default_targeting_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateVariationValue,
						},
						"environment": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateKey ,
						},
					},
				},
			},
			"default_off_targeting_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateVariationValue,
						},
						"environment": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateKey ,
						},
					},
				},
			},
			"variations": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateVariationValue,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
	return resourceImport(resourceFeatureFlagRead, d, meta)
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
	variationsKind := validateOrDefaultToBoolean(d.Get("variations_kind").(string))
	variations := d.Get("variations").([]interface{})
	defaultTargetingRule := d.Get("default_targeting_rule").([]interface{})
	defaultOffTargetingRule := d.Get("default_off_targeting_rule").([]interface{})
	customProperties := d.Get("custom_properties").([]interface{})
	transformedTags := make([]string, len(tags))
	sort.Strings(transformedTags)

	transformedVariations, err := transformVariationsFromTerraformFormat(variations, variationsKind)
	if err != nil {
		return err
	}

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
		Variations:       transformedVariations,
		CustomProperties: transformedCustomProperties,
	}

	var response JsonFeatureFlag
	err = client.Post(getFlagCreateUrl(project), payload, []int{201}, &response)
	if err != nil {
		return err
	}

	defaultTargetinRulePayload, err := createPayloadForDefaultTargeting(defaultTargetingRule, variations)
	if err != nil {
		return err
	}
	offOffTargetingRulePayload, err := createPayloadForDefaultOffTargeting(defaultOffTargetingRule, variations)
	if err != nil {
		return err
	}

	patchPayload := append(defaultTargetinRulePayload, offOffTargetingRulePayload...)

	if len(patchPayload) > 0 {
		_, err = client.Patch(getFlagUrl(project, key), patchPayload, []int{200}, NUMBER_OF_RETRY)
		if err != nil {
			return err
		}
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)
	d.Set("description", description)
	d.Set("temporary", temporary)
	d.Set("include_in_snippet", includeInSnippet)
	d.Set("tags", transformedTags)
	d.Set("variations", variations)
	d.Set("custom_properties", customProperties)
	d.Set("default_targeting_rule", defaultTargetingRule)
	d.Set("default_off_targeting_rule", defaultOffTargetingRule)

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
	transformedVariations := transformVariationsFromLaunchDarklyFormat(response.Variations)
	transformedCustomProperties := transformCustomPropertiesFromLaunchDarklyFormat(response.CustomProperties)

	d.SetId(key)
	d.Set("name", response.Name)
	d.Set("key", response.Key)
	d.Set("description", response.Description)
	d.Set("temporary", response.Temporary)
	d.Set("include_in_snippet", response.IncludeInSnippet)
	// We don't update the state if it contains the same tags as Launchdarkly (regardless of their ordering)
	if !reflect.DeepEqual(transformTagsFromTerraformFormat(d.Get("tags").([]interface{})), response.Tags) {
		d.Set("tags", response.Tags)
	}

	// This is a hack to prevent recreating a feature flag when it was created with an older
	// version of the provider that didn't supported specifying the variations kind
	if response.VariationsKind == VARIATIONS_BOOLEAN_KIND {
		d.Set("variations_kind", response.VariationsKind)
	}
	if err := d.Set("variations", transformedVariations); err != nil {
		return err
	}
	if err := d.Set("custom_properties", transformedCustomProperties); err != nil {
		return err
	}
	
	return nil
}

func resourceFeatureFlagUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(Client)
	project := resourceData.Get("project_key").(string)
	name := resourceData.Get("name").(string)
	description := resourceData.Get("description").(string)
	temporary := resourceData.Get("temporary").(bool)
	includeInSnippet := resourceData.Get("include_in_snippet").(bool)
	tags := resourceData.Get("tags").([]interface{})
	customProperties := resourceData.Get("custom_properties").([]interface{})
	variations := resourceData.Get("variations").([]interface{})
	defaultTargetingRule := resourceData.Get("default_targeting_rule").([]interface{})
	defaultOffTargetingRule := resourceData.Get("default_off_targeting_rule").([]interface{})

	transformedCustomProperties, err := transformCustomPropertiesFromTerraformFormat(customProperties)
	if err != nil {
		return err
	}

	if err := applyChangesToVariations(resourceData, client); err != nil {
		return err
	}
	
	defaultTargetingRulePayload, err := createPayloadForDefaultTargeting(defaultTargetingRule, variations)
	if err != nil {
		return err
	}
	defaultOffTargetingRulePayload, err := createPayloadForDefaultOffTargeting(defaultOffTargetingRule, variations)
	if err != nil {
		return err
	}

	targetingPayload := append(defaultTargetingRulePayload, defaultOffTargetingRulePayload...)

	mainPayload := []map[string]interface{}{{
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

	payload := append(mainPayload, targetingPayload...)

	_, err = client.Patch(getFlagUrl(project, resourceData.Id()), payload, []int{200}, NUMBER_OF_RETRY)
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
	sort.Strings(transformed)

	return transformed
}

func createPayloadForDefaultTargeting(defaultTargetingRules []interface{}, variations []interface{}) ([]map[string]interface{}, error) {
	patchPayload := make([]map[string]interface{}, len(defaultTargetingRules))
	for index, defaultTargetingRule := range defaultTargetingRules {
		targetingRule := defaultTargetingRule.(map[string]interface{})

		variationIndex, err := getDefaultVariationIndex(variations, targetingRule["value"].(string))
		if err != nil {
			return nil, err
		}

		patchPayload[index] = map[string]interface{}{
			"op":    "replace",
			"path":  fmt.Sprintf("/environments/%s/fallthrough/variation", targetingRule["environment"].(string)),
			"value": variationIndex,
		}
	}
	return patchPayload, nil
} 

func createPayloadForDefaultOffTargeting(defaultOffTargetingRules []interface{}, variations []interface{}) ([]map[string]interface{}, error) {
	patchPayload := make([]map[string]interface{}, len(defaultOffTargetingRules))
	for index, defaultOffTargetingRule := range defaultOffTargetingRules {
		targetingRule := defaultOffTargetingRule.(map[string]interface{})

		variationIndex, err := getDefaultOffVariationIndex(variations, targetingRule["value"].(string))
		if err != nil {
			return nil, err
		}
		patchPayload[index] = map[string]interface{}{
			"op":    "replace",
			"path":  fmt.Sprintf("/environments/%s/offVariation", targetingRule["environment"].(string)),
			"value": variationIndex,
		}
	}
	return patchPayload, nil
} 

func getDefaultOffVariationIndex(variations []interface{}, variationValue string) (int, error) {
	if len(variations) > 0 {
		if len(variationValue) > 0 {
			return getVariationIndex(variations, variationValue)
		}
		return len(variations) - 1, nil
	}
	return 1, nil
}

func getDefaultVariationIndex(variations []interface{}, variationValue string) (int, error) {
	if len(variations) > 0 {
		if len(variationValue) > 0 {
			return getVariationIndex(variations, variationValue)
		}
	}
	return 0, nil
}

func getVariationIndex(variations []interface{}, variationValue string) (int, error) {
	for index, rawVariationValue := range variations {
		variation := rawVariationValue.(map[string]interface{})
		if variationValue == variation[VARIATION_VALUE_KEY].(string) {
			return index, nil
		}
	}
	return len(variations), errors.New(fmt.Sprintf("%s is not a valid variation value as it is not in the provided variations", variationValue))
}

func transformVariationsFromTerraformFormat(variations []interface{}, variationsKind string) ([]JsonVariations, error) {
	transformedVariations := make([]JsonVariations, len(variations))
	for index, rawVariationValue := range variations {
		variation := rawVariationValue.(map[string]interface{})
		var value interface{}
		name := variation[VARIATION_NAME_KEY].(string)
		description := variation[VARIATION_DESCRIPTION_KEY].(string)

		if variationsKind == VARIATIONS_STRING_KIND {
			value = variation[VARIATION_VALUE_KEY].(string)
		} else if variationsKind == VARIATIONS_NUMBER_KIND {
			convertedNumberValue, err := strconv.Atoi(variation[VARIATION_VALUE_KEY].(string))
			if err != nil {
				return nil, err
			}
			value = convertedNumberValue
		} else if variationsKind == VARIATIONS_BOOLEAN_KIND {
			convertedBooleanValue, err := strconv.ParseBool(variation[VARIATION_VALUE_KEY].(string))
			if err != nil {
				return nil, err
			}
			value = convertedBooleanValue
		}

		transformedVariations[index] = JsonVariations{
			Name:        name,
			Value:       value,
			Description: description,
		}
	}

	return transformedVariations, nil
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

func applyChangesToVariations(resourceData *schema.ResourceData, client Client) error {
	project := resourceData.Get("project_key").(string)
	key := resourceData.Id()
	variations := resourceData.Get("variations").([]interface{})

	var response JsonFeatureFlag
	err := client.GetInto(getFlagUrl(project, key), []int{200}, &response)
	if err != nil {
		return err
	}
	actualNumberOfVariation := len(response.Variations)

	transformedVariations, err := transformVariationsFromTerraformFormat(variations, resourceData.Get("variations_kind").(string))
	if err != nil {
		return err
	}

	newNumberOfVariation := len(transformedVariations)

	//Remove variations
	if newNumberOfVariation < actualNumberOfVariation {
		var deletePayloadValue []interface{} = make([]interface{}, actualNumberOfVariation - newNumberOfVariation)
		for i := 0; i < len(deletePayloadValue); i++ {
			removeValue := map[string]interface{}{ 
				"op":   "remove",
				"path": fmt.Sprintf("/variations/%d", actualNumberOfVariation - 1),
			}
			deletePayloadValue[i] = removeValue
			actualNumberOfVariation--
		}
		_, err = client.Patch(getFlagUrl(project, key), deletePayloadValue, []int{200}, NUMBER_OF_RETRY)
		if err != nil {
			return err
		}
		//Update values off existing variations that wont be delete
		var updatePayloadValue []interface{} = make([]interface{}, 3*newNumberOfVariation)
		for i := 0; i < newNumberOfVariation; i++ {
			replaceValue := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/value", i),
				"value": transformedVariations[i].Value,
			}
			replaceName := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/name", i),
				"value": transformedVariations[i].Name,
			}
			replaceDescription := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/description", i),
				"value": transformedVariations[i].Description,
			}
			updatePayloadValue[i*3] = replaceValue
			updatePayloadValue[(i*3)+1] = replaceName
			updatePayloadValue[(i*3)+2] = replaceDescription
		}
		_, err = client.Patch(getFlagUrl(project, key), updatePayloadValue, []int{200}, NUMBER_OF_RETRY)
		if err != nil {
			return err
		}
	}

	//Add new variations
	if newNumberOfVariation >= actualNumberOfVariation {
		//Update values off existing variations
		var updatePayloadValue []interface{} = make([]interface{}, 3*actualNumberOfVariation)
		for i := 0; i < actualNumberOfVariation; i++ {
			replaceValue := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/value", i),
				"value": transformedVariations[i].Value,
			}
			replaceName := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/name", i),
				"value": transformedVariations[i].Name,
			}
			replaceDescription := map[string]interface{}{
				"op":    "replace",
				"path":  fmt.Sprintf("/variations/%d/description", i),
				"value": transformedVariations[i].Description,
			}
			updatePayloadValue[i*3] = replaceValue
			updatePayloadValue[(i*3)+1] = replaceName
			updatePayloadValue[(i*3)+2] = replaceDescription
		}
		_, err = client.Patch(getFlagUrl(project, key), updatePayloadValue, []int{200}, NUMBER_OF_RETRY)
		if err != nil {
			return err
		}
		
		var createPayloadValue []interface{} = make([]interface{}, newNumberOfVariation - actualNumberOfVariation)
		for i := 0; i < len(createPayloadValue); i++ {
			createPayloadValue[i] = map[string]interface{}{
				"op":    "add",
				"path":  fmt.Sprintf("/variations/%d", actualNumberOfVariation + i),
				"value": transformedVariations[actualNumberOfVariation + i],
			}
		}
		_, err = client.Patch(getFlagUrl(project, key), createPayloadValue, []int{200}, NUMBER_OF_RETRY)
		if err != nil {
			return err
		}
	}

	return nil
}

func transformVariationsFromLaunchDarklyFormat(properties []JsonVariations) interface{} {
	transformedVariations := make([]map[string]interface{}, 0)

	for _, variation := range properties {
		transformedVariation := make(map[string]interface{})
		transformedVariation["name"] = variation.Name
		transformedVariation["value"] = fmt.Sprint(variation.Value)
		transformedVariation["description"] = variation.Description

		transformedVariations = append(transformedVariations, transformedVariation)
	}

	return transformedVariations
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

func validateOrDefaultToBoolean(variationsKind string) string {
	if len(variationsKind) > 0 {
		return variationsKind
	}

	return DEFAULT_VARIATIONS_KIND
}
