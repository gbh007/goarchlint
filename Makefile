CMD = $(shell go env GOBIN)/goarchlint

ifeq ($(shell go env GOBIN),)
CMD := $(shell go env GOPATH)/bin/goarchlint
endif

.PHONY: docself
docself:
	go run cmd/goarchlint/main.go generate -o doc/arch

.PHONY: install
install:
	go build -o $(CMD) cmd/goarchlint/main.go