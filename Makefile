MODELVERSION ?= "latest"

# Run tests.
test: fmt lint tidy test-cmd
			$(TEST_COMMAND)
			# Nested go modules are not tested for some reason, so test them separately
			cd sdk && $(TEST_COMMAND)

test-cmd:
ifeq (, $(shell which gotestsum))
TEST_COMMAND=go test ./...
else
TEST_COMMAND=gotestsum
endif

tidy:
	go mod tidy
	cd sdk && go mod tidy

# Bundle static assets
generate:
	go run asset_generator.go

build:
	go build -o bin/riser

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
lint:
	golangci-lint run
	cd sdk && golangci-lint run

# compile and run unit tests on change.
# requires filewatcher and gotestsum
watch:
	filewatcher gotestsum

# updates to the latest api models
# Note: As of go 1.13 GOSUMDB returns a 410. Disabling until we figure out why.
update-model:
	GOSUMDB=off go get -u github.com/riser-platform/riser-server/api/v1/model@$(MODELVERSION) && go mod tidy
	cd sdk && GOSUMDB=off go get -u github.com/riser-platform/riser-server/api/v1/model@$(MODELVERSION) && go mod tidy

	# Github actions passes the full ref so strip it off
VERSIONCLEAN=$(subst refs/tags/,,$(VERSION))
release: check-version
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s -X 'main.versionString=$(VERSIONCLEAN)'" -o="bin/darwin-amd64/riser"
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X 'main.versionString=$(VERSIONCLEAN)'" -o="bin/linux-amd64/riser"
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -s -X 'main.versionString=$(VERSIONCLEAN)'" -o="bin/windows-amd64/riser.exe"
	zip -r riser-darwin-amd64.zip -j bin/darwin-amd64/riser
	zip -r riser-linux-amd64.zip -j bin/linux-amd64/riser
	zip -r riser-windows-amd64.zip -j bin/windows-amd64/riser.exe


check-version:
	@if test -z "${VERSION}"; then echo "Usage: make <target> VERSION=<version>"; exit 1; fi

# Warning! This deletes and recreates the minikube project named "demo"!
# Useful for testing demo installation from scratch
minikube-demo: build
	minikube delete -p demo
	minikube start -p demo
	riser demo install

# Uses the riser binary in your path. Requires both the riser context and kube context to be configured to the target environment
e2e:
	go test -count=1 -v ./pkg/e2e

