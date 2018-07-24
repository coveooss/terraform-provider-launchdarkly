VERSION = 1.0.0
SOURCES = $(wildcard *.go)

.PHONY: default
default: build cross-compile

.PHONY: build
build:
	go build -o terraform-provider-launchdarkly_v$(VERSION)

.PHONY: clean
clean:
	rm terraform-provider-launchdarkly_v*

.PHONY: cross-compile
cross-compile:
	GOOS=windows GOARCH=amd64 go build -o terraform-provider-launchdarkly_v$(VERSION)-windows
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-launchdarkly_v$(VERSION)-osx
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-launchdarkly_v$(VERSION)-linux

.PHONY: test
test: build
	terraform init
	terraform apply
	terraform destroy
