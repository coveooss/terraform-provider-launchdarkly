SOURCES = $(wildcard *.go)

.PHONY: build
build:
	go build -o terraform-provider-launchdarkly

clean:
	rm terraform-provider-launchdarkly

.PHONY: test
test: build
	terraform init
	terraform apply
	terraform destroy
