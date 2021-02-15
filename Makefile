all: build

# tools
TOOLS_MOD_DIR := ./tools
TOOLS_DIR := $(abspath ./.tools)
$(TOOLS_DIR)/golangci-lint: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	@echo BUILD golangci-lint
	@cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/mockery: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	@echo BUILD mockery
	@cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/mockery github.com/vektra/mockery/v2

# lint, test
.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint
	@echo LINT
	@$(TOOLS_DIR)/golangci-lint run -c .github/linters/.golangci.yaml --out-format colored-line-number
	@printf "LINT... \033[0;32m [OK] \033[0m"

.PHONY: test
test: gen
	@echo SMALL TEST
	@go test -v -test.short ./...
	@printf "SMALL TEST... \033[0;32m [OK] \033[0m"

test-medium: gen
	@echo MEDIUM TEST
	@rm -rf test/report
	@go test -v ./...
	@printf "MEDIUM TEST... \033[0;32m [OK] \033[0m"

.PHONY: test-with-coverage
test-with-coverage:
	@rm -rf test/apitest
	@go test -v -test.short ./... -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt -o coverage.html

# build
BIN_DIR := $(abspath ./bin)
BUILD_TARGETS=build/server
build: $(BUILD_TARGETS)

.PHONY: $(BIN_DIR)
$(BIN_DIR):
	@mkdir -p $@

.PHONY: $(BUILD_TARGETS)
build/server: $(BIN_DIR)
	@echo BUILD server
	@go build -v -o ./bin/server ./cmd/server

# gen
GEN_TARGETS=gen/mock
.PHONY: gen
gen: $(GEN_TARGETS)

.PHONY: $(GEN_TARGETS)
gen/mock: $(TOOLS_DIR)/mockery
	@echo GENERATE mocks
	@find ./usecase -type d -name mocks | xargs rm -rf
	@go generate ./...

