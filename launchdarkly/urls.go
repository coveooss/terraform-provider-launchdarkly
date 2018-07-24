package launchdarkly

import "fmt"

const rootUrl = "https://app.launchdarkly.com/api/v2"

func getProjectCreateUrl() string {
	return fmt.Sprintf("%s/projects", rootUrl)
}

func getProjectUrl(project string) string {
	return fmt.Sprintf("%s/projects/%s", rootUrl, project)
}

func getFlagCreateUrl(project string) string {
	return fmt.Sprintf("%s/flags/%s", rootUrl, project)
}

func getFlagUrl(project string, flag string) string {
	return fmt.Sprintf("%s/flags/%s/%s", rootUrl, project, flag)
}

func getEnvironmentCreateUrl(project string) string {
	return fmt.Sprintf("%s/projects/%s/environments", rootUrl, project)
}

func getEnvironmentUrl(project string, environment string) string {
	return fmt.Sprintf("%s/projects/%s/environments/%s", rootUrl, project, environment)
}
