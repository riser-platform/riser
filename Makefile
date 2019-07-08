build:
	go build -o bin/riser

test: fmt lint
	gotestsum

# Runs teh server in development mode
run: fmt lint
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
lint:
	golangci-lint run

# compile and run unit tests on change.
# requires filewatcher and gotestsum
watch:
	filewatcher gotestsum

# updates to the latest api models
update-model:
	go get -u github.com/tshak/riser-server/api/v1/model
	go mod tidy



