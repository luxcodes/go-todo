# go-todo

A maybe not so simple Go-based TODO REST API

---

## Features

- CRUD operations for TODO items
- PostgreSQL database with migrations
- RESTful API with JSON
- Auto-generated Swagger (OpenAPI) docs
- Interactive Swagger UI (served via Nginx, using CDN)
- Integration tests using Testcontainers
- Docker Compose support for local development
- Nginx reverse proxy for unified API and docs access

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.24+
- [Docker](https://www.docker.com/get-started)
- [Make](https://www.gnu.org/software/make/)
- [swag](https://github.com/swaggo/swag) for Swagger docs (`go install github.com/swaggo/swag/cmd/swag@latest`)

---

## Building & Running

Start all services using Docker Compose:

```sh
make docker-up
```

- API endpoints: [http://localhost/](http://localhost/)
- Swagger UI: [http://localhost/swagger/](http://localhost/swagger/)

---

## API Documentation

- **Swagger UI** is served via Nginx at `/swagger/index.html`.
- The UI loads the OpenAPI JSON from `/swagger/swagger.json`.
- To regenerate Swagger docs after editing handler comments:

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
docs/              # Swagger docs and UI (swagger.json, swagger.html, etc.)
.env               # Environment variables
Makefile           # Build and test commands
docker-compose.yml # Docker Compose config
nginx.conf         # Nginx reverse proxy config
```

---

## License

MIT
