SOURCES = $(wildcard *.go)
TEST?=./...

.PHONY: default
default: build cross-compile

.PHONY: build
build:
	go build
	go test $(TEST) -timeout=30s -parallel=4

.PHONY: clean
clean:
	rm -f terraform-provider-launchdarkly
	rm -rf output

.PHONY: cross-compile
cross-compile:
	GOOS=windows GOARCH=amd64 go build -o output/windows_amd64/terraform-provider-launchdarkly
	tar -C output/windows_amd64 -czf output/terraform-provider-launchdarkly_windows_amd64.tar.gz terraform-provider-launchdarkly
	GOOS=darwin GOARCH=amd64 go build -o output/osx_amd64/terraform-provider-launchdarkly
	tar -C output/osx_amd64 -czf output/terraform-provider-launchdarkly_osx_amd64.tar.gz terraform-provider-launchdarkly
	GOOS=linux GOARCH=amd64 go build -o output/linux_amd64/terraform-provider-launchdarkly
	tar -C output/linux_amd64 -czf output/terraform-provider-launchdarkly_linux_amd64.tar.gz terraform-provider-launchdarkly

.PHONY: test
test:
	go test $(TEST) -timeout=30s -parallel=4
