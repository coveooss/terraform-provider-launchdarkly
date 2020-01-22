package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKey,
			},
		},
	}
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	d.SetId(d.Id())
	d.Set("key", d.Id())

	resourceProjectRead(d, meta)

	return []*schema.ResourceData{d}, nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	name := d.Get("name").(string)
	key := d.Get("key").(string)

	payload := JsonProject{
		Name: name,
		Key:  key,
	}

	var response JsonProject
	err := client.Post(getProjectCreateUrl(), payload, []int{201}, &response)
	if err != nil {
		return err
	}

	// Default environments will be created, we want to get rid of those
	environmentKeys, err := getEnvironmentKeys(client, key)
	if err != nil {
		return err
	}

	err = ensureThereIsADummyEnvironment(client, key)
	if err != nil {
		return err
	}

	for _, environmentKey := range environmentKeys {
		err = client.Delete(getEnvironmentUrl(key, environmentKey), []int{204})
		if err != nil {
			return err
		}
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)

	client := m.(Client)

	var response JsonProject
	err := client.GetInto(getProjectUrl(key), []int{200}, &response)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(key)
	d.Set("name", response.Name)
	d.Set("key", response.Key)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	name := d.Get("name").(string)

	payload := []map[string]string{{
		"op":    "replace",
		"path":  "/name",
		"value": name,
	}}

	_, err := client.Patch(getProjectUrl(d.Id()), payload, []int{200}, 0)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	err := client.Delete(getProjectUrl(d.Id()), []int{204, 404})
	if err != nil {
		return err
	}

	return nil
}
