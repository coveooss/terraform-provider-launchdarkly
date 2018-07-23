package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,

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
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	color := d.Get("color").(string)

	payload := map[string]string{
		"name":  name,
		"key":   key,
		"color": color,
	}

	err := client.Post("/projects/"+project+"/environments", payload, map[int]bool{201: true}, nil)
	if err != nil {
		return err
	}

	// If a dummy environment was created before, we no longer need it
	err = ensureThereIsNoDummyEnvironment(client, project)
	if err != nil {
		return err
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)
	d.Set("color", color)

	return nil
}

func resourceEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project_key").(string)
	key := d.Get("key").(string)

	client := m.(Client)

	payload := make(map[string]interface{})

	err := client.Get("/projects/"+project+"/environments/"+key, map[int]bool{200: true}, &payload)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", payload["name"])
	d.Set("key", payload["key"])
	d.Set("color", payload["color"])

	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)
	name := d.Get("name").(string)
	color := d.Get("color").(string)

	payload := []map[string]string{{
		"op":    "replace",
		"path":  "/name",
		"value": name,
	}, {
		"op":    "replace",
		"path":  "/color",
		"value": color,
	}}

	err := client.Patch("/projects/"+project+"/environments/"+d.Id(), payload, map[int]bool{200: true}, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	project := d.Get("project_key").(string)

	onlyOne, err := isThereOnlyOneEnvironment(client, project)
	if err != nil {
		return err
	}

	if onlyOne {
		println("Creating dummy environment since we cannot delete the last environment in a project")

		err = ensureThereIsADummyEnvironment(client, project)
		if err != nil {
			return err
		}
	}

	err = client.Delete("/projects/"+project+"/environments/"+d.Id(), map[int]bool{204: true, 404: true})
	if err != nil {
		return err
	}

	return nil
}
