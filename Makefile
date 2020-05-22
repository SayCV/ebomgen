# Makefile for the `libasciidoc` project

# tools
CUR_DIR=$(shell pwd)
INSTALL_PREFIX=$(CUR_DIR)/bin
VENDOR_DIR=vendor
SOURCE_DIR ?= .
COVERPKGS := $(shell go list ./... | grep -v vendor | paste -sd "," -)

DEVTOOLS=

ifeq ($(OS),Windows_NT)
BINARY_PATH=$(INSTALL_PREFIX)/ebomgen.exe
else
BINARY_PATH=$(INSTALL_PREFIX)/ebomgen
endif

# Call this function with $(call log-info,"Your message")
define log-info =
@echo "INFO: $(1)"
endef


.PHONY: help
# Based on https://gist.github.com/rcmachado/af3db315e31383502660
## Display this help text.
help:/
	$(info Available targets)
	$(info -----------------)
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		helpCommand = substr($$1, 0, index($$1, ":")-1); \
		if (helpMessage) { \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			gsub(/##/, "\n                                     ", helpMessage); \
		} else { \
			helpMessage = "(No documentation)"; \
		} \
		printf "%-35s - %s\n", helpCommand, helpMessage; \
		lastLine = "" \
	} \
	{ hasComment = match(lastLine, /^## (.*)/); \
          if(hasComment) { \
            lastLine=lastLine$$0; \
	  } \
          else { \
	    lastLine = $$0 \
          } \
        }' $(MAKEFILE_LIST)

.PHONY: install-devtools
## Install development tools.
install-devtools:
	@go mod download

$(INSTALL_PREFIX):
# Build artifacts dir
	@mkdir -p $(INSTALL_PREFIX)

.PHONY: prebuild-checks
## Check that all tools where found
prebuild-checks: $(INSTALL_PREFIX)

.PHONY: generate
## generate the .go file based on the asciidoc grammar
generate: prebuild-checks
	@echo "generating the parser..."

.PHONY: generate-optimized
## generate the .go file based on the asciidoc grammar
generate-optimized:
	@echo "generating the parser (optimized)..."

.PHONY: test
## run all tests excluding fixtures and vendored packages
test: generate-optimized
	@echo $(COVERPKGS)

.PHONY: test-with-coverage
## run all tests excluding fixtures and vendored packages
test-with-coverage: generate-optimized
	@echo $(COVERPKGS)

.PHONY: test-fixtures
## run all fixtures tests
test-fixtures: generate-optimized
	@echo "run all fixtures tests"

.PHONY: bench-parser
## run the benchmarks on the parser
bench-parser: generate-optimized
	@echo "run the benchmarks on the parser"

.PHONY: build
## build the binary executable from CLI
build: $(INSTALL_PREFIX) generate-optimized
	$(eval BUILD_COMMIT:=$(shell git rev-parse --short HEAD))
	$(eval BUILD_TAG:=$(shell git tag --contains $(BUILD_COMMIT)))
	$(eval BUILD_TIME:=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ'))
	@echo "building $(BINARY_PATH) (commit:$(BUILD_COMMIT) / tag:$(BUILD_TAG) / time:$(BUILD_TIME))"
	@go build -ldflags \
	  " -X github.com/saycv/ebomgen.BuildCommit=$(BUILD_COMMIT)\
	    -X github.com/saycv/ebomgen.BuildTag=$(BUILD_TAG) \
	    -X github.com/saycv/ebomgen.BuildTime=$(BUILD_TIME)" \
	  -o $(BINARY_PATH) \
	  cmd/ebomgen/*.go

.PHONY: lint
## run golangci-lint against project
lint:
	@golangci-lint run -E gofmt,golint,megacheck,misspell ./...

PARSER_DIFF_STATUS :=

.PHONY: verify-parser
## verify that the parser was built with the latest version of pigeon, using the `optimize-grammar` option
verify-parser: prebuild-checks
ifneq ($(shell git diff --quiet pkg/parser/parser.go; echo $$?), 0)
	$(error "parser was generated with an older version of 'mna/pigeon' or without the '-optimize' option(s).")
else
	@echo "parser is ok"
endif

.PHONY: install
## installs the binary executable in the $GOPATH/bin directory
install: install-devtools build
	@cp $(BINARY_PATH) $(GOPATH)/bin

.PHONY: quick-install
## installs the binary executable in the $GOPATH/bin directory without prior tools setup and parser generation
quick-install:
	$(eval BUILD_COMMIT:=$(shell git rev-parse --short HEAD))
	$(eval BUILD_TAG:=$(shell git tag --contains $(BUILD_COMMIT)))
	$(eval BUILD_TIME:=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ'))
	@echo "building $(BINARY_PATH) (commit:$(BUILD_COMMIT) / tag:$(BUILD_TAG) / time:$(BUILD_TIME))"
	@go build -ldflags \
	  " -X github.com/saycv/ebomgen.BuildCommit=$(BUILD_COMMIT)\
	    -X github.com/saycv/ebomgen.BuildTag=$(BUILD_TAG) \
	    -X github.com/saycv/ebomgen.BuildTime=$(BUILD_TIME)" \
	  -o $(BINARY_PATH) \
	  cmd/ebomgen/*.go
	@cp $(BINARY_PATH) $(GOPATH)/bin
