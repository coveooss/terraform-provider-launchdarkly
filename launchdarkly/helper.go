package launchdarkly

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func parseCompositeID(id string) (p1 string, p2 string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		p1 = parts[0]
		p2 = parts[1]
	} else {
		err = fmt.Errorf("error: Import composite ID requires two parts separated by colon, eg x:y")
	}
	return
}

type importFunc func(d *schema.ResourceData, meta interface{}) error

func resourceImport(readMethod importFunc, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectKey, resourceKey, err := parseCompositeID(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(resourceKey)
	d.Set("project_key", projectKey)
	d.Set("key", resourceKey)

	readMethod(d, meta)

	return []*schema.ResourceData{d}, nil
}
