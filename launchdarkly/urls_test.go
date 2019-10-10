package launchdarkly

import (
	"testing"
)

func TestGetProjectCreateUrl(t *testing.T) {
	expectedUrl := "https://app.launchdarkly.com/api/v2/projects"
	returnedUrl := getProjectCreateUrl()
	if returnedUrl != expectedUrl {
		t.Errorf("getProjectCreateUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}

func TestGetProjectUrl(t *testing.T) {
	aProjectName := "my-project"
	expectedUrl := "https://app.launchdarkly.com/api/v2/projects/" + aProjectName
	returnedUrl := getProjectUrl(aProjectName)
	if returnedUrl != expectedUrl {
		t.Errorf("getProjectUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}

func TestGetFlagCreateUrl(t *testing.T) {
	aProjectName := "my-project"
	expectedUrl := "https://app.launchdarkly.com/api/v2/flags/" + aProjectName
	returnedUrl := getFlagCreateUrl(aProjectName)
	if returnedUrl != expectedUrl {
		t.Errorf("getFlagCreateUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}

func TestGetFlagUrl(t *testing.T) {
	aProjectName := "my-project"
	aFlagName := "my-super-flag"
	expectedUrl := "https://app.launchdarkly.com/api/v2/flags/" + aProjectName + "/" + aFlagName
	returnedUrl := getFlagUrl(aProjectName, aFlagName)
	if returnedUrl != expectedUrl {
		t.Errorf("getFlagUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}

func TestGetEnvironmentCreateUrl(t *testing.T) {
	aProjectName := "my-project"
	expectedUrl := "https://app.launchdarkly.com/api/v2/projects/" + aProjectName + "/environments"
	returnedUrl := getEnvironmentCreateUrl(aProjectName)
	if returnedUrl != expectedUrl {
		t.Errorf("getEnvironmentCreateUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}

func TestGetEnvironmentUrl(t *testing.T) {
	aProjectName := "my-project"
	anEnvironmentName := "my-marvelous-environment"
	expectedUrl := "https://app.launchdarkly.com/api/v2/projects/" + aProjectName + "/environments/" + anEnvironmentName
	returnedUrl := getEnvironmentUrl(aProjectName, anEnvironmentName)
	if returnedUrl != expectedUrl {
		t.Errorf("getEnvironmentUrl expected return value was '%s' but got '%s'", expectedUrl, returnedUrl)
	}
}