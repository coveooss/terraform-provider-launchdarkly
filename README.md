[![Build Status](https://travis-ci.com/coveo/terraform-provider-launchdarkly.svg?branch=master)](https://travis-ci.com/coveo/terraform-provider-launchdarkly)
[![license](http://img.shields.io/badge/license-Apache-brightgreen.svg)](https://github.com/coveo/terraform-provider-launchdarkly/blob/master/LICENSE)

# A Terraform provider for LaunchDarkly feature flags 

This provider allows creating Projects, Environments and Feature Flags in LaunchDarkly.

## Getting Started

### Installation

Download the appropriate binary from the GitHub release, and install it on your local computer as described [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

### Sample usage

Have a look at the `main.tf` file for a sample configuration using the provider.

#### Importing resources
Using the command `import` you need to follow the following syntax.

For resources `environment` and `feature_flag` :
You need 2 values in the resource import ID separated by `:` . 

The project key in the resource key.
e.g.: `import launchdarkly_environment.my-env critical-updates-dev:dev`

For the `project` resource you only need the project key. e.g.: `import launchdarkly_project.my-project critical-updates-dev`

## Building the provider
Clone the repository, and run `make` at the root of the working copy.
