# Targets:
#  - build:  Build 'exec-assert' binary
#  - lint:   Run linters
#  - check:  Run unit tests
#  - test:   Run all tests
#  - verify: Run all tests and linters
#  - clean:  Clean up
.PHONY: build lint check test verify clean

# Build 'exec-assert' binary
# 
# Args:
#   GOFLAGS: Extra flags to pass to 'go build'
#
# Examples:
#   make build
#   make build GOFLAGS=-v
build:
	go build $(GOFLAGS) .

# Run linters
#
# Args:
#   GOVETFLAGS: Extra flags to pass to 'go vet'
#
# Examples:
#   make lint
#   make link GOVETFLAGS=-v
#   make link GOVETFLAGS=-shadowstrict
lint:
	go tool vet -all -shadow $(GOVETFLAGS) pkg

# Run unit tests
#
# Args:
#   GOFLAGS: Extra flags to pass to 'go test'
#
# Examples:
#   make check
#   make check GOFLAGS=-cover
check:
	go test $(GOFLAGS) -v ./...

# Run all tests
#
# Args:
#   GOTESTFLAGS: Extra flags to pass to 'go test'
#
# Examples:
#   make test
#   make test GOTESTFLAGS=-cover
test:
	$(MAKE) check GOFLAGS=$(GOTESTFLAGS)
	test/cmd.sh

# Run all tsts and linters
#
# Args:
#   GOTESTFLAGS: Extra flags to pass to 'go test'
#   GOVETFLAGS: Extra flags to pass to 'go vet'
#
# Examples:
#   make verify
#   make verify GOTESTFLAGS=-cover
#   make verify GOVETFLAGS=-shadowstrict
verify:
	$(MAKE) lint GOVETFLAGS=$(GOVETFLAGS)
	$(MAKE) test GOTESTFLAGS=$(GOTESTFLAGS)
.PHONY: verify

# Clean up all build artifacts
#
# Example:
#   make clean
clean:
	rm -f exec-assert
