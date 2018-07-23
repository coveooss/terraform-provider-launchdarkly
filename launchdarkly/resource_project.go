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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	name := d.Get("name").(string)
	key := d.Get("key").(string)

	payload := map[string]string{
		"name": name,
		"key":  key,
	}

	err := client.Post("/projects", payload, map[int]bool{201: true}, nil)
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
		err = client.Delete("/projects/"+key+"/environments/"+environmentKey, map[int]bool{204: true})
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

	payload := make(map[string]interface{})

	err := client.Get("/projects/"+key, map[int]bool{200: true}, &payload)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", payload["name"])
	d.Set("key", payload["key"])

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

	err := client.Patch("/projects/"+d.Id(), payload, map[int]bool{200: true}, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	err := client.Delete("/projects/"+d.Id(), map[int]bool{204: true, 404: true})
	if err != nil {
		return err
	}

	return nil
}
