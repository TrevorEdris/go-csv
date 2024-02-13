SHELL ?= /bin/bash
export GIT_HASH =$(shell git rev-parse --short HEAD)
export VERSION ?=  $(shell printf "`./.tools/version`${VERSION_SUFIX}")
export DEV_DOCKER_COMPOSE ?= docker-compose.dev.yaml

all: build

.PHONY: help
help: ## List of available commands
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\033[36m\1\\033[m:\2/' | column -c2 -t -s :)"

# --------------------------[ Development ]--------------------------
.PHONY: up
up: ## Run the application and follow the logs
	docker-compose -f ${DEV_DOCKER_COMPOSE} up -d
	make -s logs

.PHONY: down
down: ## Stop all containers
	docker-compose -f ${DEV_DOCKER_COMPOSE} down

.PHONY: restart
restart: down up ## Restart all containers

.PHONY: logs
logs: ## Print logs in stdout
	docker-compose -f ${DEV_DOCKER_COMPOSE} logs -f app

# -----------------------------[ Build ]-----------------------------

.PHONY: version
version: ## Automatically calculate the semantic version based on the number of commits since the last change to the VERSION file
	@echo ${VERSION}

.PHONY: build
build:
	mkdir -p bin
	go build -ldflags "-s -w ${LDFLAGS} -X main.Version=${VERSION}" -o bin/gocsv main.go

# ----------------------------[ Install ]----------------------------

.PHONY: install
install:
	go install -ldflags="-X main.Version=${VERSION}" ./cmd/gocsv
