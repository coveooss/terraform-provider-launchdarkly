provider "launchdarkly" {
  access_token = "api-3963b908-ec5a-4d9e-adb3-6ee30a119724"
}

resource "launchdarkly_project" "my-project" {
  key = "my-project-key"
  name = "My Super Project!!"
}