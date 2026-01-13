SHELL = /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin

LINTER_NAME := golangci-lint
LINTER_VERSION := v2.7.2

.PHONY: all build test debezium-compose-up debezium-compose-down load-test-data vendor install-linter lint fmt tools tools-update generate activate-python-venv install-admin-deps run-django-admin

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/bot/main.go
	go build -mod=vendor -o $(BIN_PATH)/manager ./cmd/manager/main.go
	go build -mod=vendor -o $(BIN_PATH)/msgconsumer ./cmd/msgconsumer/main.go

test:
	go test ./...

debezium-compose-up:
	docker compose -f ./script/docker/debezium-docker-compose.yml up --build

debezium-compose-down:
	docker compose -f ./script/docker/debezium-docker-compose.yml down

load-test-data:
	docker run -it --rm --network host \
		-v ./script/db/dataset/test_data.sql:/script/test_data.sql \
		alpine/psql:17.7 \
		"postgresql://bot:bot@localhost:5432/bot" -f /script/test_data.sql

vendor:
	go mod tidy
	go mod vendor

install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml

fmt: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) fmt --config .golangci.yml

tools: install-linter
	@if [ ! -f $(GOBIN)/mockgen ]; then\
		echo "Installing mockgen";\
		GOBIN=$(GOBIN) go install go.uber.org/mock/mockgen@v0.6.0;\
	fi

tools-update:
	go get tool

generate:
	$(ENV_PATH) go generate ./...

activate-python-venv:
	@if [ ! -d $(BIN_PATH)/python_venv ]; then \
		python -m venv $(BIN_PATH)/python_venv; \
	fi

install-admin-deps: activate-python-venv
	source $(BIN_PATH)/python_venv/bin/activate && \
		python -m pip install \
		Django==5.2 \
		python-decouple==3.8 \
		psycopg==3.1.18 \
		psycopg2-binary

run-django-admin: install-admin-deps
	source $(BIN_PATH)/python_venv/bin/activate && \
		python cmd/adminpanel/manage.py migrate && \
		python cmd/adminpanel/manage.py runserver \
