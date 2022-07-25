#!make
MAKEFLAGS += --silent

COMMIT_SHA=$(shell git rev-parse --short HEAD)
DOCKER_REGISTRY=registry.digitalocean.com/mmess
IMAGE_NAME=rental

.PHONY: show_version setup migrate seed_db api_v1_gen test deploy build_image push_image do_registry_login all


show_version:
	@echo ${COMMIT_SHA}

setup: # setup local development environment
	go install github.com/rakyll/gotest@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

test: # run tests
	docker-compose up -d database
	sleep 2 # wait for database to be ready, TODO: find a way to make this deterministic
	gotest -v ./...
	docker-compose down -v

deploy: # Deploy the app through Helm
	DOCKER_REGISTRY=${DOCKER_REGISTRY} \
	COMMIT_SHA=${COMMIT_SHA} \
	./deployment/bin/deploy.sh

build_image:
	docker build --no-cache --platform linux/amd64 -t ${IMAGE_NAME}:${COMMIT_SHA} .
	docker tag ${IMAGE_NAME}:${COMMIT_SHA} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${COMMIT_SHA}

do_registry_login:
	 doctl registry login --expiry-seconds 600

push_image: do_registry_login
	docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${COMMIT_SHA}

seed_db: migrate
	cat db/seed/data.sql | docker-compose exec -T database psql "postgres://rental:rental@localhost:5432/rental?sslmode=disable" -

migrate:
	docker-compose up -d database
	sleep 2 # wait for database to be ready, TODO: find a way to make this deterministic
	migrate -path db/migrations/ -database "postgres://rental:rental@localhost:5432/rental?sslmode=disable" up

api_v1_gen:
	oapi-codegen -package gen api/rental-v1.0.yml > pkg/api/gen/api.gen.go

all: build_image push_image
