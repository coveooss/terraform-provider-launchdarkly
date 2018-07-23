SOURCES = $(wildcard *.go)

.PHONY: build
build: terraform-provider-launchdarkly

clean:
	rm terraform-provider-launchdarkly

.PHONY: test
test: build
	terraform init
	terraform apply
	terraform destroy -parallelism=1

terraform-provider-launchdarkly: $(SOURCES)
	go build -o terraform-provider-launchdarkly
