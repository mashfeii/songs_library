COVERAGE_FILE ?= coverage.out

TARGET ?= songs_library

.PHONY: build
build:
	@echo "Выполняется go build для таргета ${TARGET}"
	@mkdir -p bin
	@go build -o ./bin/${TARGET} ./cmd/${TARGET}

.PHONY: docs
docs:
	@swag init -g ./cmd/songs_library/main.go -o docs/

## test: run all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/mashfeii/songs_library/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'
	@go tool cover -html='$(COVERAGE_FILE)' -o coverage.html && xdg-open coverage.html

