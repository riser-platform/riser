# Run tests.
test: fmt lint test-cmd
			$(TEST_COMMAND)
			# Nested go modules are not tested for some reason, so test them separately
			cd sdk && $(TEST_COMMAND)

test-cmd:
ifeq (, $(shell which gotestsum))
TEST_COMMAND=go test ./...
else
TEST_COMMAND=gotestsum
endif

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
	GOSUMDB=off go get -u github.com/tshak/riser-server/api/v1/model@latest
	go mod tidy
	cd sdk && GOSUMDB=off go get -u github.com/tshak/riser-server/api/v1/model@latest && go mod tidy



