# go-todo

A simple Go-based TODO REST API with Postgres, Swagger documentation, and integration tests.

---

## Features

- CRUD operations for TODO items
- PostgreSQL database with migrations
- RESTful API with JSON
- Auto-generated Swagger (OpenAPI) docs
- Integration tests using Testcontainers
- Docker Compose support for local development

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.24+
- [Docker](https://www.docker.com/get-started)
- [Make](https://www.gnu.org/software/make/)
- [swag](https://github.com/swaggo/swag) for Swagger docs (`go install github.com/swaggo/swag/cmd/swag@latest`)

### Database Setup

Start Postgres w/ migrations using Docker Compose:

```sh
make docker-up
```

---

## Building & Running

### Build the application

```sh
make build
```

### Run the application

```sh
make run
```

The API will be available at [http://localhost:8080](http://localhost:8080).

---

## API Documentation

Swagger UI is available at:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

To regenerate Swagger docs after editing handler comments:

```sh
swag init -g cmd/api/main.go
```

---

## Testing

### Unit & Integration Tests

Run all tests:

```sh
make test
```

Run only integration tests (database):

```sh
make itest
```

---

## Useful Make Commands

| Command         | Description                      |
|-----------------|----------------------------------|
| `make build`    | Build the application            |
| `make run`      | Run the application              |
| `make test`     | Run all tests                    |
| `make itest`    | Run integration tests            |
| `make docker-up`| Start DB container               |
| `make docker-down` | Stop DB container             |
| `make clean`    | Remove built binaries            |
| `make watch`    | Live reload with Air (if installed) |

---

## Project Structure

```
cmd/api/           # Main application entrypoint
internal/          # Application code (server, database, models)
migrations/        # SQL migration files
docs/              # Swagger docs (auto-generated)
.env               # Environment variables
Makefile           # Build and test commands
docker-compose.yml # Docker Compose config
```

---

## License

MIT
