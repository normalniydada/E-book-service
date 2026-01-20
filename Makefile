.PHONY: test cover build clean help

BINARY_NAME=beauty-salon
COVERAGE_FILE=coverage.out

test:
	go test -v ./...

cover:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)
	@echo "Opening HTML report..."
	go tool cover -html=$(COVERAGE_FILE)

race:
	go test -race ./...

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build-secure

build-secure:
	go install mvdan.cc/garble@latest
	garble -literals -tiny -seed=random build -o bin/beauty-salon-secure ./cmd/main.go