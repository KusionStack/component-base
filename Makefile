# Copyright 2024 KusionStack Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash

# Define variables so `make --warn-undefined-variables` works.
PRINT_HELP ?=

define VET_HELP_INFO
# Run 'go vet' command to vet Go code.
#
# Example:
#   make vet
endef
.PHONY: vet
ifeq ($(PRINT_HELP),y)
vet:
	@echo "$$VET_HELP_INFO"
else
vet:
	go vet ./...
endif

define FMT_HELP_INFO
# Run 'go fmt' command to format Go code.
#
# Example:
#   make fmt
endef
.PHONY: fmt
ifeq ($(PRINT_HELP),y)
fmt:
	@echo "$$FMT_HELP_INFO"
else
fmt:
	go fmt ./...
endif

define TEST_HELP_INFO
# Build and run tests.
#
# Example:
#   make test
endef
.PHONY: test
ifeq ($(PRINT_HELP),y)
test:
	@echo "$$TEST_HELP_INFO"
else
test: fmt vet
	go test ./... -coverprofile coverage.out
endif