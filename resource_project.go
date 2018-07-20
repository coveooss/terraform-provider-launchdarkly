package main

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

	payload := map[string]string {
		"name" : name,
		"key": key,
	}

	err := client.Post("/projects", payload, 201, nil)
	if err != nil {
		return err
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

	err := client.Get("/projects/" + key, 200, &payload)
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

	payload := []map[string]string { {
		"op": "replace",
		"path": "/name",
		"value": name,
	} }

	err := client.Patch("/projects/" + d.Id(), payload, 200, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	err := client.Delete("/projects/" + d.Id(), 204)
	if err != nil {
		return err
	}

	return nil
}
