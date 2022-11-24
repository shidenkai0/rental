# Rental Service API (showcase)

## Introduction

Rental is an API for handling car rental services. It is a showcase of an e2e implementation of a simple API in Golang, and is not intended to be useful or performant as is, but rather as an example of how to implement an API with Golang, from backend to infrastructure to documentation. It includes a simple Makefile, local run through docker-compose, API documentation (with codegen from the doc), infrastructure-as-code defined with Terraform, CI/CD deployments through Github Actions. It allows the user to create, update, and delete cars, as well as customers. It also allows the user to rent and return cars.

## API documentation

A full OpenAPI v3/Swagger specification can be found locally at `api/rental-v1.0.yaml`.
You can use a tool such as Swagger or Postman to interact with the API when running it locally at http://localhost:9090 or in production at https://rental.mmess.dev (as specified in the API Spec).

## Configuration
The following environment variables are available for configuration:

* `PORT`: The port to run the API on. Defaults to `9090`.
* `DATABASE_URL`: The URL of the database to use. Defaults to `postgres://rental:rental@localhost:5432/rental`.
* `DATABASE_MAX_OPEN_CONNS`: The maximum number of open connections to the database. Defaults to 5.
* `DATABASE_MAX_IDLE_CONNS`: The maximum number of idle connections to the database. Defaults to 2.
* `BASIC_AUTH_USER`: The username to use for basic authentication. Defaults to `rental`.
* `BASIC_AUTH_PASSWORD`: The password to use for basic authentication. Defaults to `rental`.

## Developing locally

### Prerequisites

Development on this project assumes that you have Docker and docker-compose setup, as well as a functional Go 1.18 or greater development environment, to set these up:

- [Docker and docker-compose](https://docs.docker.com/compose/install/)
- [Go](https://go.dev/doc/install)

Before starting development on this project, install the recommended development tools by running:

```bash
make setup
```
### Running the API

Proceed with setting up the local database by running:

```
docker-compose up -d database
make seed_db
```

Then, you can run the following commands to start the API:

```bash
docker-compose up
```

The API Server should now be listening on http://localhost:9090, the v1 API is available under the path `/v1`.

The default basic auth credentials when running locally are `rental:rental`.

### Running the tests
To run all the unit tests, run:

```bash
make test
```

## Adding a new feature

### Without API breaking changes
- Start by editing the API Spec at `api/rental-v1.0.yaml`
- Generate API boilerplate code by running:
```bash
make api_v1_gen
```
- If you are just changing the logic of an endpoint, simply edit the appropriate method in `pkg/api.Server`
- If you added some endpoints, implement them on the `pkg/api.Server` type, use the `pkg/api/gen.ServerInterface` interface as a reference for which methods to implement.
### With API breaking changes
Introducing breaking changes into the API requires creating a new version of the API Spec. This can be done by copying the existing API Spec and modifying it, for instance:

```bash
cp api/rental-v1.0.yaml api/rental-v2.0.yaml
```
Then, add a new command to the Makefile named openapi_v${MAJOR_VERSION}_gen using the `api_v1_gen` command as a reference, and changing the version number where necessary.

All that is left is then to implement a Server for this new version of the API in `pkg/api` using version 1 of the API as a reference.

Do not forget to register the API under a new group in the Echo Server:
```go
e.Group("/v2")
```
This will allow the new version of the API to be accessible under the `/v2` path.


## Deploying to production

Deploying to production is done automatically by the CI, however, to deploy locally, and assuming you have access to the production environment, you can run:

```bash
make deploy
```

## Database migrations

### Create a migration

To create a migration called `$migration_name` run the following command:

```
migrate create -ext sql -dir db/migrations -seq $migration_name
```

This should create files called `$migration_name.up.sql` and `$migration_name.down.sql` in the `db/migrations` directory. Fill in the necessary SQL statements to create the tables and columns you need. Do not forget to fill the "down" migration with the SQL statements to drop the tables and columns you created: running the "up" and "down" files sequentially should result in a no-op.

### Run migrations
#### Locally

```
make migrate
```

#### In production
Migrations are run automatically by the CI for now, however if you still want to run them manually, you can still do so by following the steps below:


- Make sure the current kubernetes context is set to the production cluster
- Run the following command:

```
make migrate_prod
```

If you wish to use more advanced [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) commands, for instance to set the database to a specific revision, you can use the following command, customizing it yo your needs:

```
kubectl run migration -it --restart=Never --image ${DOCKER_REGISTRY}/${MIGRATION_IMAGE_NAME}:${COMMIT_SHA} --rm -- -database "${DATABASE_URL}" $YOUR_CUSTOM_FLAGS_AND_ARGUMENTS
```
