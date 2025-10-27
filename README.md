# Shorten URL Service# Shorten URL Service

Clean Architecture Go Project TemplateA clean Go web service template following best practices.

## Project Structure## Quick Start

````bash

project/# 1. Copy environment variables

├── cmd/                         # Command-related filescp .env.example .env

│   └── app/                     # Application entry point

│       └── main.go              # Main application logic# 2. Install dependencies

├── internal/                    # Internal codebasego mod tidy

│   ├── delivery/                # External layer (HTTP Handlers, gRPC, etc.)

│   │   └── http/                # HTTP delivery# 3. Run the application

│   │       └── user_handler.go  # User-specific HTTP handlergo run main.go

│   ├── usecases/                # Use Cases (business logic layer)```

│   │   └── user_service.go      # User-specific service logic

│   ├── repository/              # Repository (data access, external services)## Project Structure

│   │   └── user_repo.go         # User-specific data access

│   └── entities/                # Entities (core models, domain objects)```

│       └── user.go              # User model├── cmd/                 # Application entrypoints

├── pkg/                         # Shared utilities or helpers├── internal/            # Private application code

├── configs/                     # Configuration files│   ├── app/            # Application composition root

├── go.mod                       # Go module definition│   ├── config/         # Configuration

└── go.sum                       # Go module checksum file│   ├── handler/        # HTTP handlers

```│   ├── middleware/     # HTTP middlewares

│   ├── model/          # Data models

## Clean Architecture Layers│   ├── repository/     # Data access layer

│   └── service/        # Business logic

1. **Entities** (`internal/entities/`) - Core business models├── pkg/                # Public libraries

2. **Use Cases** (`internal/usecases/`) - Business logic└── web/                # Web assets (optional)

3. **Repository** (`internal/repository/`) - Data access layer```

4. **Delivery** (`internal/delivery/`) - External interfaces (HTTP, gRPC, etc.)

## API Endpoints

## Getting Started

- `GET /health` - Health check

```bash- `GET /api/v1/ping` - API status

cd cmd/app

go run main.go## Development

```
```bash
# Run tests
go test ./...

# Build
go build -o bin/app main.go

# Format code
go fmt ./...
```

## Environment Variables

See `.env.example` for available configuration options.
````
