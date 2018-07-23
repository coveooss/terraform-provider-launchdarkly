package launchdarkly

import "strconv"

const dummyEnvironmentKey = "dummy-environment"

func isThereOnlyOneEnvironment(client Client, project string) (bool, error) {
	payload := make(map[string]interface{})

	err := client.Get("/projects/"+project, map[int]bool{200: true}, &payload)
	if err != nil {
		return false, err
	}

	environments := payload["environments"].([]interface{})

	println("There are currently " + strconv.Itoa(len(environments)) + " environments in project " + project)

	return len(environments) == 1, nil
}

func getEnvironmentKeys(client Client, project string) ([]string, error) {
	payload := make(map[string]interface{})

	err := client.Get("/projects/"+project, map[int]bool{200: true}, &payload)
	if err != nil {
		return nil, err
	}

	environments := payload["environments"].([]interface{})

	keys := make([]string, len(environments))

	for i, env := range environments {
		keys[i] = env.(map[string]interface{})["key"].(string)
	}

	return keys, nil
}

func ensureThereIsADummyEnvironment(client Client, project string) error {
	exists, err := isThereADummyEnvironment(client, project)
	if err != nil {
		return err
	}

	if !exists {
		return createDummyEnvironment(client, project)
	} else {
		return nil
	}
}

func ensureThereIsNoDummyEnvironment(client Client, project string) error {
	exists, err := isThereADummyEnvironment(client, project)
	if err != nil {
		return err
	}

	if exists {
		println("A dummy environment was found, deleting it")
		return deleteDummyEnvironment(client, project)
	} else {
		println("No dummy environment was found")
		return nil
	}
}

func isThereADummyEnvironment(client Client, project string) (bool, error) {
	statusCode, err := client.GetStatus("/projects/" + project + "/environments/" + dummyEnvironmentKey)
	if err != nil {
		return false, err
	}

	return statusCode == 200, nil
}

func createDummyEnvironment(client Client, project string) error {
	println("Creating dummy environment")

	payload := map[string]string{
		"name":  dummyEnvironmentKey,
		"key":   dummyEnvironmentKey,
		"color": "FFFFFF",
	}

	err := client.Post("/projects/"+project+"/environments", payload, map[int]bool{201: true}, nil)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyEnvironment(client Client, project string) error {
	println("Deleting the dummy environment")

	err := client.Delete("/projects/"+project+"/environments/"+dummyEnvironmentKey, map[int]bool{204: true, 404: true})
	if err != nil {
		return err
	}

	return nil
}
