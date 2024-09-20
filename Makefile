
POSTGRES_APP_USER ?= praktikum
POSTGRES_APP_PASS ?= pass
POSTGRES_APP_DB ?= goph_keeper

DATABASE_URL ?= postgresql://${POSTGRES_APP_USER}:${POSTGRES_APP_PASS}@localhost:5432/${POSTGRES_APP_DB}
#######################################################################################################################

DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)
#######################################################################################################################

## ▸▸▸ Development commands ◂◂◂

.PHONY: help
help:			## Show this help
	@fgrep -h "## " $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/## //'

.PHONY: clean
clean:			## Remove generated artifacts
	@rm -rf ./bin
	@rm -rf ./docker/volume

#######################################################################################################################

## ▸▸▸ Docker commands ◂◂◂

.PHONY: config
config:			## Show Docker config
	docker compose ${DOCKER_COMPOSE_FILES} config

.PHONY: up
up:			## Run Docker services
	docker compose ${DOCKER_COMPOSE_FILES} up --detach

.PHONY: down
down:			## Stop Docker services
	docker compose ${DOCKER_COMPOSE_FILES} down

.PHONY: ps
ps:			## Show Docker containers info
	docker ps --size --all --filter "name=url-shortener-api"

#######################################################################################################################

## ▸▸▸ Utils commands ◂◂◂

.PHONY: connect
connect:		## Connect to the database
	pgcli ${DATABASE_URL}

.PHONY: goose-status
goose-status:		## Dump the migration status for the current DB
	goose -dir ./server/migrations postgres ${DATABASE_URL} status

PHONY: goose-up
goose-up:		## Migrate the DB to the most recent version available
	goose -dir ./server/migrations postgres ${DATABASE_URL} up

PHONY: goose-down
goose-down:		## Roll back the version by 1
	goose -dir ./server/migrations postgres ${DATABASE_URL} down

test: build-with-coverage
	@rm -fr .coverdata
	@mkdir -p .coverdata
	@go test ./...  -coverpkg=./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool covdata percent -i=.coverdata

#check-coverage: test
#	@go tool covdata textfmt -i=.coverdata -o profile.txt
#	@go tool cover -html=profile.txt
#
build:
	@go build -C client -o client

build-with-coverage:
	@go build -C client -cover -o client

.DEFAULT_GOAL := build
#######################################################################################################################
