package launchdarkly

import (
	"strconv"
)

const dummyEnvironmentKey = "dummy-environment"

func getEnvironmentKeys(client Client, project string) ([]string, error) {
	var response JsonProject
	err := client.GetInto(getProjectUrl(project), []int{200}, &response)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, env := range response.Environments {
		keys = append(keys, env.Key)
	}

	return keys, nil
}

func ensureWeCanDeleteEnvironment(client Client, project string) error {
	onlyOne, err := isThereOnlyOneEnvironment(client, project)
	if err != nil {
		return err
	}

	if onlyOne {
		println("Creating dummy environment since we cannot delete the last environment in a project")
		return ensureThereIsADummyEnvironment(client, project)
	} else {
		return nil
	}
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
	statusCode, err := client.GetStatus(getEnvironmentUrl(project, dummyEnvironmentKey))
	if err != nil {
		return false, err
	}

	return statusCode == 200, nil
}

func isThereOnlyOneEnvironment(client Client, project string) (bool, error) {
	var response JsonProject
	err := client.GetInto(getProjectUrl(project), []int{200}, &response)
	if err != nil {
		return false, err
	}

	println("There are currently " + strconv.Itoa(len(response.Environments)) + " environments in project " + project)

	return len(response.Environments) == 1, nil
}

func createDummyEnvironment(client Client, project string) error {
	println("Creating dummy environment")

	payload := JsonEnvironment{
		Name:  dummyEnvironmentKey,
		Key:   dummyEnvironmentKey,
		Color: "FFFFFF",
	}

	var response JsonEnvironment
	err := client.Post(getEnvironmentCreateUrl(project), payload, []int{201}, &response)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyEnvironment(client Client, project string) error {
	println("Deleting the dummy environment")

	err := client.Delete(getEnvironmentUrl(project, dummyEnvironmentKey), []int{204, 404})
	if err != nil {
		return err
	}

	return nil
}
