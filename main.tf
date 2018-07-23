provider "launchdarkly" {
  access_token = "api-3963b908-ec5a-4d9e-adb3-6ee30a119724"
}

resource "launchdarkly_project" "my-project" {
  key = "my-project-key"
  name = "test"
}

resource "launchdarkly_environment" "dev" {
  project_key = "${launchdarkly_project.my-project.key}"
  name = "Development"
  key = "dev"
  color = "FF0000"
}

resource "launchdarkly_environment" "hipaa" {
  project_key = "${launchdarkly_project.my-project.key}"
  name = "HIPAA"
  key = "hipaa"
  color = "FF00FF"
}

resource "launchdarkly_feature_flag" "my-flag" {
  project_key = "${launchdarkly_project.my-project.key}"
  key = "my-flag"
  name = "My Super Flag"
  description = "description!!"
  tags = ["foo", "bar", "spam"]
  custom_properties = {
    "some.property" = "blah"
  }
}
