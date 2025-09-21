CMD = $(shell go env GOBIN)/goarchlint

ifeq ($(shell go env GOBIN),)
CMD := $(shell go env GOPATH)/bin/goarchlint
endif

.PHONY: selfdoc
selfdoc:
	go run cmd/goarchlint/main.go generate

.PHONY: selflint
selflint:
	go run cmd/goarchlint/main.go run

.PHONY: install
install:
	go build -o $(CMD) cmd/goarchlint/main.go