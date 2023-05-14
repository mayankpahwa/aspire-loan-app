BUILD_DIR := $(CURDIR)/build
BINARY ?= $(BUILD_DIR)/aspire-loan-app

.PHONY: build
build:  ## Build service image. Needs APP=[api|consumer] env variable
	docker build -f Dockerfile -t aspire-loan-app .

.PHONY: build-binary
build-binary:  ## Build service binary
	go build -o $(BINARY) ./cmd/api
