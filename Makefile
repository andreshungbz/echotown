# Project Makefile

# default rule run when no arguments are provided to make
.DEFAULT_GOAL := build

# declare following rules to not be associated with any real files
.PHONY: format vet build clean

# variables
PREFIX = [make]
BUILD_DIR = bin/
LOG_DIR = log/

# rules

clean:
	@echo "$(PREFIX) removing $(BUILD_DIR) and log directory"
	@rm -rfv $(BUILD_DIR)
	@rm -rfv $(LOG_DIR)

format: clean
	@echo "$(PREFIX) formatting source code"
	@go fmt ./...

vet: format
	@echo "$(PREFIX) vetting source code"
	@go vet ./...

build: vet
	@echo "$(PREFIX) building executable in $(BUILD_DIR)"
	@go build -o $(BUILD_DIR)
