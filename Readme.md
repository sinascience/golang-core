# Venturo Golang Core

A robust, slim, and well-organized backend codebase core using Golang. This project is built with a "Fat Model" (Active Record) pattern for simplicity and speed, containerized with Docker, and includes a dedicated migration tool.

## Table of Contents

  - [Codebase Structure](https://www.google.com/search?q=%23codebase-structure)
  - [Getting Started](https://www.google.com/search?q=%23getting-started)
      - [Prerequisites](https://www.google.com/search?q=%23prerequisites)
      - [Configuration](https://www.google.com/search?q=%23configuration)
      - [Running with Docker (Recommended)](https://www.google.com/search?q=%23running-with-docker-recommended)
      - [Running Locally](https://www.google.com/search?q=%23running-locally)
  - [Database Migrations](https://www.google.com/search?q=%23database-migrations)
  - [Creating a New Feature](https://www.google.com/search?q=%23creating-a-new-feature)
  - [Concurrency (Goroutines)](https://www.google.com/search?q=%23concurrency-goroutines)
  - [Core Libraries](https://www.google.com/search?q=%23core-libraries)

-----

## Codebase Structure

The project follows a clean architecture pattern designed for simplicity and scalability.

```
/
├── cmd/                  # Application entry points (main packages)
│   ├── migrate/          # The database migration tool
│   └── server/           # The main API server
├── configs/              # Configuration loading (.env)
├── database/
│   └── migrations/       # SQL migration files
├── internal/
│   ├── database/         # Database connection and migration logic
│   ├── handler/http/     # HTTP Handlers (Controllers)
│   ├── model/            # Data models with database logic (Fat Model/Active Record)
│   ├── server/           # Server setup, dependency injection, and routing
│   └── service/          # Business logic services
├── pkg/
│   └── utils/            # Shared utility functions (JWT, hashing, etc.)
├── .air.toml             # Configuration for hot-reloading
├── docker-compose.yml    # Docker service definitions
└── Dockerfile            # Docker build instructions
```

-----

## Getting Started

### Prerequisites

  * Go 1.24+
  * Docker & Docker Compose
  * MySQL (only if running locally without Docker)

### Configuration

1.  Copy the example environment file: `cp .env.example .env`
2.  Edit the `.env` file with your specific configuration.

### Running with Docker (Recommended)

This is the easiest way to get started. It handles both the Go application and the MySQL database.

1.  **Build and Run the Services:**

    ```bash
    docker-compose up --build
    ```

    For hot-reloading during development, the setup is already configured. Changes to `.go` files will automatically restart the server inside the container.

2.  **Run Migrations:** Open a new terminal and run:

    ```bash
    docker-compose exec app migrate up
    ```

The API server will be available at `http://localhost:3000`.

### Running Locally

1.  **Ensure MySQL is Running:** Make sure your local MySQL server is active.
2.  **Update `.env`:** Set `DB_HOST=127.0.0.1` or your local IP.
3.  **Run Migrations:**
    ```bash
    go run cmd/migrate/main.go up
    ```
4.  **Run the Server:**
    ```bash
    go run cmd/server/main.go
    ```

-----

## Database Migrations

The project uses a dedicated Go program to handle migrations.

  * **Apply all `up` migrations:**
    ```bash
    # Docker
    docker-compose exec app migrate up

    # Local
    go run cmd/migrate/main.go up
    ```
  * **Roll back the last migration:**
    ```bash
    # Docker
    docker-compose exec app migrate down

    # Local
    go run cmd/migrate/main.go down
    ```
  * **Drop all tables and re-apply all migrations:**
    ```bash
    # Docker
    docker-compose exec app migrate fresh

    # Local
    go run cmd/migrate/main.go fresh
    ```

-----

## Creating a New Feature

Here is the step-by-step pattern for adding a new feature (e.g., "Posts").

1.  **Create the Model (`internal/model/post_model.go`):**
    Define the `Post` struct and add its standard CRUD methods (`Save`, `FindAll`, etc.), following the pattern in `user_model.go`.

2.  **Create the Migration:**
    Generate the migration files:

    ```bash
    # This command is run locally
    migrate create -ext sql -dir database/migrations -seq create_posts_table
    ```

    Edit the generated `.up.sql` and `.down.sql` files with the schema for your `posts` table.

3.  **Create the Service (`internal/service/post_service.go`):**
    Create a service file to contain any business logic related to posts.

4.  **Create the Handler (`internal/handler/http/post_handler.go`):**
    Create the handler to manage HTTP requests and responses for posts. It will call the `PostService`.

5.  **Register the Routes (`internal/server/routes.go`):**
    In the `registerRoutes` function, initialize your new service and handler, then define the endpoints for your feature.

    ```go
    // In registerRoutes function...
    postService := service.NewPostService(db)
    postHandler := http.NewPostHandler(postService)

    api.Get("/posts", postHandler.GetAll)
    api.Post("/posts", postHandler.Create)
    ```

-----

## Concurrency (Goroutines)

This section will be updated with guidelines for implementing concurrent tasks. The general approach will be to use goroutines for background tasks (e.g., sending emails, processing images) initiated from the service layer. Channels will be used for communication and synchronization where necessary, potentially using a `sync.WaitGroup` to manage groups of goroutines.

-----

## Core Libraries

  * [**Fiber**](https://github.com/gofiber/fiber): Web framework.
  * [**GORM**](https://github.com/go-gorm/gorm): Database ORM.
  * [**golang-migrate**](https://github.com/golang-migrate/migrate): Database migrations.
  * [**go-jwt**](https://github.com/golang-jwt/jwt): JSON Web Token implementation.
  * [**Bcrypt**](https://pkg.go.dev/golang.org/x/crypto/bcrypt): Password hashing.
  * [**Air**](https://github.com/cosmtrek/air): Live-reloading for development.
  * [**godotenv**](https://github.com/joho/godotenv): `.env` file loading.