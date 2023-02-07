PROJECT_NAME := "fred"
PKG := "git.tu-berlin.de/mcc-fred/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /ext/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v /ext/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint megalint container docs

all: build

lint: ## Lint the files
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed, check https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint -E revive --timeout 5m0s run

megalint: ## Megalint all files
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed, check https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint -E asciicheck -E depguard -E dogsled -E dupl -E errorlint -E exhaustive -E exportloopref -E forbidigo -E funlen  -E gochecknoinits -E gocognit -E gocritic -E gocyclo -E gofmt -E revive -E gomnd -E gomodguard -E goprintffuncname -E gosec -E interfacer -E makezero -E maligned -E misspell -E nestif -E nlreturn -E stylecheck -E unconvert -E unparam run

test: ## Run unittests
	@rm -rf pkg/badgerdb/test.db
	@go test -short ${PKG_LIST}
	@rm -rf pkg/badgerdb/test.db

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@sh ./ci/tools/coverage.sh;

dep: ## Get the dependencies
	@go mod download

build: dep ## Build the binary file
	@go build -v $(PKG)/cmd/frednode

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

container: ## Create a Docker container
	@docker build . -t git.tu-berlin.de:5000/mcc-fred/fred/fred

docs: ## Build the FogStore documentation
	@mdpdf docs/doc.md

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
