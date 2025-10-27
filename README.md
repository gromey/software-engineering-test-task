# Simple CRUD Interface

Rewrite the README according to the application.

The task itself can be found [here](/TASK.md)

## Prerequisites

- [Docker](https://www.docker.com/get-started/)
- [Goose](https://github.com/pressly/goose)
- [Gosec](https://github.com/securego/gosec)

## Getting Started

1. Start database

Install the latest version of Docker Compose:

```shell
sudo apt-get install -y docker-compose-v2
```

```
## Via Makefile
make db

## Via Docker
docker compose up -d db
```

2. Run migrations

```
## Via Makefile
make migrate-up

## Via Goose
DB_DRIVER=postgres
DB_STRING="host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
goose -dir ./migrations $(DB_DRIVER) $(DB_STRING) up
```

3. Run application

```shell
cp .env.example .env
```

Set variables if necessary

> Note:
> - The API uses an X-API-Key header for authentication.
> - By default, the API key is set to `secret`.
> - Include this header in every request, for example:
> ```shell
> curl -H "X-API-Key: secret" http://localhost:8080/api/v1/users
> ```

```
## Via Go
go run cmd/main.go

## Via Makefile
make run

## Via Makefile in docker
make up
```