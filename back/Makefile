PREFIX ?= $(shell pwd)
SOURCE_FILES ?= ./...

TEST_PATTERN ?= .
TEST_OPTS ?=

.PHONY: setup
setup: ## Install dev tools
	@brew install golangci/tap/golangci-lint
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/psampaz/go-mod-outdated
	@go get github.com/swaggo/swag/cmd/swag

.PHONY: test
test: ## Run all the tests
	@go test $(TEST_OPTS) -v -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@go tool cover -html=coverage.out

.PHONY: tidy
tidy: ## Prune any no-longer-needed dependencies
	@go mod tidy

.PHONY: fmt
fmt: tidy ## Run gofmt and goimports on all go files
	@find . -name '*.go' | while read -r file; do gofmt -w -s "$${file}"; goimports -w "$${file}"; done

.PHONY: lint
lint: fmt ## Run all the linters
	golangci-lint run --enable-all

.PHONY: swagger
swagger: ## Generate swagger definition
	@swag init --parseInternal

.PHONY: build
build: swagger ## Build a test binary
	@go build

.PHONY: release
release: swagger ## Build a release binary for distribution
	@GOOS=linux GOARCH=amd64 go build

.PHONY: deploy
deploy: clean fmt release ## Deploy a release binary
	@scp back keepthemoment:
	@ssh keepthemoment 'sudo systemctl stop back'
	@ssh keepthemoment 'sudo rm /home/www/keepthemoment/back/back'
	@ssh keepthemoment 'sudo mv back /home/www/keepthemoment/back/back'
	@ssh keepthemoment 'sudo chown nginx:nginx /home/www/keepthemoment/back/back'
	@ssh keepthemoment 'sudo systemctl start back'

.PHONY: clean
clean: ## Remove all generated files
	@rm -rf coverage.out back
	@go clean

.PHONY: stat
stat: fmt ## Count number of code ines into each file
	find . -name '*.go' | xargs wc -l

.PHONY: check
check: ## Checks whether direct dependency updates are avalible
	go list -u -m -json all | "${GOPATH}/bin/go-mod-outdated" -update -direct

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
default: help