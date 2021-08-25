# ===================================================
# Makefile for building pxe-init service
# ===================================================

# ----------------
# Configuration
# ----------------

DOCKER  := docker
GO      := go
PROTOC  := protoc

# use vendored dependencies
export GOFLAGS     := -mod=vendor

# enable buildkit for better build features
export DOCKER_BUILDKIT         := 1

DIST_DIR     := ./dist
CMD_DIR      := ./cmd
INTERNAL_DIR := ./internal
HELM_DIR     := ./helm

TEST_MODULES               = $(shell $(GO) list $(INTERNAL_DIR)/...)
INTERNAL_NON_TEST_GO_FILES = $(shell find $(INTERNAL_DIR) -type f -name '*.go' -not -name '*_test.go')
INTERNAL_GO_FILES          = $(shell find $(INTERNAL_DIR) -type f -name '*.go')
CMD_GO_FILES               = $(shell find $(CMD_DIR) -type f -name '*.go')

VERSION = $(shell (git describe --long --tags --match 'v[0-9]*' || echo v0.0.0) | cut -c2-)
LDFLAGS = -X main.version=$(VERSION)

DOCKER_IMAGE := pxe-init-server
DOCKER_TAG = $(VERSION)

PUSH_REGISTRY    := zbox.mortcloud.com:5000

PUSH_PROJECT     := pxe-init
PUSH_IMAGE       := $(DOCKER_IMAGE)
PUSH_TAG          = $(DOCKER_TAG)


# oci annotations to add to the image as labels
LABEL_CREATED   = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LABEL_AUTHORS  := thomasmortensson@gmail.com
LABEL_SOURCE   := https://github.com/thomasmortensson/pxe-init
LABEL_VERSION   = $(VERSION)
LABEL_REVISION  = $(COMMIT)

# tracking for built images
BUILT := .built

# the user and group for running the image
USER  = $(shell id -u)
GROUP = $(shell id -g)

# the command to use for running the image by default
DOCKER_CMD := sh

MOCKERY      := GOFLAGS="" mockery
MOCKERY_ARGS := --all --keeptree --dir $(INTERNAL_DIR)
MOCKS_DIR    := ./mocks

# helper functions
built          = $(DOCKER) images $(1) --format '{{.ID}}' >> $(BUILT)
push           = $(DOCKER) push $(1)

# ----------------
# Binary/App Targets
# ----------------

## help: Output this message and exit.
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'
.PHONY: help

# all: lint, build
.PHONY: all
all: lint build

## build: build the project
.PHONY: build
build:
	@mkdir -p dist
	$(GO) build -ldflags "$(LDFLAGS)" -o dist/ ./cmd/...

## tidy: tidy dependencies
.PHONY: tidy
tidy:
	$(GO) mod tidy

## vendor: download vendored dependencies
.PHONY: vendor
vendor:
	$(GO) mod vendor

## fmt: format the code
.PHONY: fmt
fmt:
	gci -w $(CMD_GO_FILES) $(INTERNAL_GO_FILES)
	goimports -w $(CMD_GO_FILES) $(INTERNAL_GO_FILES)

## clean: clean up built code
.PHONY: clean
clean:
	rm -rf $(DIST_DIR)
	rm -f c.out
	rm -f unit-test-output.txt
	rm -f unit_test_coverage.html

## trivy: scan for vulnerabilities
.PHONY: trivy
trivy:
	trivy fs --exit-code 0 --severity UNKNOWN,LOW,MEDIUM --no-progress --skip-dirs tests .
	trivy fs --exit-code 1 --severity HIGH,CRITICAL --no-progress --skip-dirs tests .

## lint: run the project linters
.PHONY: lint
lint:
	golangci-lint run

## mocks: Generate the mocks for all internal interfaces for testing
mocks: $(INTERNAL_NON_TEST_GO_FILES)
	rm -rf $(MOCKS_DIR)_maketemp/
	# Mockery returns error code 0 on these errors but produces incorrect output
	if $(MOCKERY) $(MOCKERY_ARGS) --output $(MOCKS_DIR)_maketemp 2>&1 | grep ERR; then \
		rm -rf $(MOCKS_DIR)_maketemp; \
		exit 1; \
	fi
	rm -rf $(MOCKS_DIR)/
	mv $(MOCKS_DIR)_maketemp $(MOCKS_DIR)

## test: Run unit tests
.PHONY: test
test: 
	$(GO) test $(TEST_MODULES) \
		-v \
		-cover \
		-coverprofile=c.out \
		-count=1 \
		| tee unit-test-output.txt
	go tool cover -html=c.out -o unit_test_coverage.html
	go tool cover -func c.out

## proto: Generate grpc code from proto
.PHONY: proto
proto:
	$(PROTOC) \
		--go_out=internal/adapters/grpc/v1 \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/adapters/grpc/v1 \
		--go-grpc_opt=paths=source_relative \
		pxe_init.proto

# ----------------
# Docker Targets
# ----------------

## docker-build: Build the pxe-init-server image.
.PHONY: docker-build
docker-build:
	$(DOCKER) build \
		--pull \
		--force-rm \
		--target $(DOCKER_IMAGE) \
		--label org.opencontainers.image.created=$(LABEL_CREATED) \
		--label org.opencontainers.image.authors=$(LABEL_AUTHORS) \
		--label org.opencontainers.image.source=$(LABEL_SOURCE) \
		--label org.opencontainers.image.version=$(LABEL_VERSION) \
		--label org.opencontainers.image.revision=$(LABEL_REVISION) \
		--label org.opencontainers.image.title=$(DOCKER_IMAGE) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		.

## docker-push: Push the docker image.
.PHONY: docker-push
docker-push:
	$(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(PUSH_REGISTRY)/$(PUSH_PROJECT)/$(PUSH_IMAGE):$(PUSH_TAG)
	$(call push,$(PUSH_REGISTRY)/$(PUSH_PROJECT)/$(PUSH_IMAGE):$(PUSH_TAG))