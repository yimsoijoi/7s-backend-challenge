# 7solutions Backend Challenge

# User Service (Go + Hexagonal Architecture)

This project is a simple **User Management API** written in **Go**.
It demonstrates **Hexagonal / Clean Architecture**, MongoDB persistence, JWT authentication, REST and gRPC adapters, and proper unit testing.

The goal of this project is **clarity and correctness**, not framework magic.

---

## Features

* User registration and login
* JWT authentication (HS256)
* CRUD operations for users
* REST API (HTTP)
* gRPC API (*TODO*)
* MongoDB persistence (official Go driver)
* Background goroutine to count users
* Unit tests with mocked interfaces
* Docker and Docker Compose support

---

## Project Structure

```
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── adapters
│   │   ├── grpc
│   │   │   ├── server.go
│   │   │   └── user.proto
│   │   ├── http
│   │   │   ├── handler.go
│   │   │   └── middleware.go
│   │   └── mongo
│   │       ├── user_document.go
│   │       ├── user_repository_test.go
│   │       └── user_repository.go
│   ├── application
│   │   ├── user_service_test.go
│   │   └── user_service.go
│   ├── domain
│   │   └── user.go
│   ├── infrastructure
│   │   ├── jwt_test.go
│   │   ├── jwt.go
│   │   ├── mocks
│   │   │   └── jwt.go
│   │   └── mongo.go
│   └── ports
│       ├── mocks
│       │   └── user_repository.go
│       ├── repository.go
│       └── service.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### Layer Responsibilities

* **Domain**
  Pure business models (e.g. `User`). No frameworks, no databases.

* **Ports**
  Interfaces that define what the core needs and provides.

  * Repository interfaces (outbound)
  * Service interfaces (inbound)

* **Application**
  Implements business use cases. Orchestrates domain logic.

* **Adapters**
  External interfaces (HTTP, gRPC, MongoDB).

* **Infrastructure**
  Technical details like JWT and MongoDB client setup.

---

## Architecture Principles

* Dependency Inversion (core does not depend on infrastructure)
* Interfaces owned by the core
* Infrastructure implements interfaces
* Easy to test (no DB or JWT in unit tests)

---

## Requirements

* Go 1.20+
* Docker & Docker Compose

---

## Configuration

All configuration is provided via **environment variables** (defined in `docker-compose.yml`):

* `APP_PORT` – HTTP server port
* `MONGO_URI` – MongoDB connection string
* `MONGO_DB` – MongoDB database name
* `JWT_SECRET` – JWT signing secret
* `JWT_TTL` – JWT expiration duration

---

## Run with Docker

```bash
docker compose up --build
```

Services:

* API: [http://localhost:8080](http://localhost:8080)
* MongoDB: localhost:27017

---

## REST API Endpoints

### Register

```
POST /register
```

```json
{
  "name": "John",
  "email": "john@test.com",
  "password": "secret"
}
```

---

### Login

```
POST /login
```

```json
{
  "email": "john@test.com",
  "password": "secret"
}
```

Response:

```json
{ "token": "<jwt>" }
```

---

### Protected Endpoints

Add header:

```
Authorization: Bearer <jwt>
```

* `GET /users`
* `GET /users/{id}`
* `PUT /users/{id}`
* `DELETE /users/{id}`

---

## gRPC

* Proto file: `internal/adapters/grpc/user.proto`
* gRPC server runs in the same application
* (Optional) JWT can be passed via metadata

---

## Background Task

A goroutine runs every **10 seconds** and logs the number of users in MongoDB.

---

## Testing

Run all tests:

```bash
go test ./...
```

### Testing Strategy

* Application layer tests mock ports
* MongoDB adapter tests use `mtest`
* JWT tested independently

No real database is required for unit tests.

---

## MongoDB

MongoDB collections are created automatically on first use. No manual migration is required.

To inspect data:

```bash
docker exec -it mongo mongosh
use app
show collections
db.users.find().pretty()
```

---

## Why This Architecture?

* Easy to reason about
* Easy to test
* Easy to replace MongoDB, HTTP, or JWT
* Suitable for real production systems

---

## Notes

This project is intentionally simple and explicit.
No ORM, no frameworks, no hidden magic.

---

## License

MIT
