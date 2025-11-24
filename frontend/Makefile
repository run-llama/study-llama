BIN="./tmp/bin"
BIN_NAME="frontend"
MAIN_PKG="."
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: test lint format format-check typecheck build format-fe build-fe clean-build-fe test-fe lint-fe install-deps-fe

format-fe:
	$(info ******************** checking formatting for frontend ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint-fe:
	$(info ******************** running lint tools for frontend ********************)
	golangci-lint run -v

test-fe: install-deps-fe
	$(info ****************** running tests for frontend ******************)
	go test ./...

install-deps-fe:
	$(info ******************** downloading dependencies for frontend ********************)
	go get -v ./...

build-fe: install-deps-fe
	$(info ****************** building frontend ******************)
	@mkdir -p ${BIN}
	GOARCH=amd64 GOOS=linux go build -o ${BIN}/${BIN_NAME}-linux-amd64 ${MAIN_PKG}
	GOARCH=amd64 GOOS=windows go build -o ${BIN}/${BIN_NAME}-windows-amd64.exe ${MAIN_PKG}
	GOARCH=arm64 GOOS=darwin go build -o ${BIN}/${BIN_NAME}-darwin-arm64 ${MAIN_PKG}

clean-build-fe:
	@rm -rf ${BIN}