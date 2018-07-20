SOURCES = $(wildcard *.go)

.PHONY: build
build: terraform-provider-launchdarkly

clean:
	rm terraform-provider-launchdarkly

.PHONY: test
test: build
	go test ./...

terraform-provider-launchdarkly: $(SOURCES)
	go build -o terraform-provider-launchdarkly
