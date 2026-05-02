# Backend Structure

This backend is a Go API service. The project is organized so startup code, HTTP routing, configuration, database access, and data models each have a clear place.

## Directory Overview

```text
backend/
  cmd/
    api/
      main.go
  internal/
    config/
    database/
    handlers/
    middleware/
    models/
    repository/
    routes/
  migrations/
  scripts/
```

## `cmd/api/main.go`

`main.go` is the application entry point. When the backend starts, execution begins here.

It should contain startup wiring only:

- Load configuration.
- Connect to the database.
- Create the Gin router.
- Register routes.
- Start the HTTP server.

It should not contain:

- SQL queries.
- Request validation logic.
- Business rules.
- Large route lists.
- Model definitions.

The goal is for `main.go` to stay small and easy to scan.

## `internal`

`internal` contains application-private Go packages. Go gives this folder special meaning: packages inside `internal` cannot be imported freely by code outside this backend module.

Use `internal` for code that belongs to this backend application, not for public reusable libraries.

## `internal/config`

The config package reads environment variables and turns them into a typed `Config` struct.

It is responsible for:

- Loading local environment values.
- Reading required settings such as database connection values.
- Applying safe defaults such as the local server port.
- Returning clear errors when required settings are missing.

It should not:

- Connect to the database.
- Start the server.
- Know about HTTP routes.
- Contain business logic.

Other packages should avoid calling `os.Getenv` directly unless there is a strong reason. Prefer adding config values here first.

## `internal/database`

The database package creates and verifies the database connection.

It is responsible for:

- Creating the PostgreSQL connection pool.
- Pinging the database during startup.
- Returning the database pool to the application.

It should not:

- Contain SQL queries for specific features.
- Know about users, platforms, roles, or other domain concepts.
- Know about HTTP requests or Gin.

Feature-specific database queries should eventually live in a repository package, not in `internal/database`.

## `internal/routes`

The routes package maps HTTP paths and methods to handlers.

Example responsibility:

```go
router.GET("/live", handlers.Live())
router.GET("/ready", handlers.Ready(db))
```

It is responsible for:

- Defining API endpoint paths.
- Grouping related endpoints.
- Applying route-level middleware.
- Connecting routes to handler functions.

It should not:

- Parse request bodies.
- Run SQL queries.
- Contain business logic.
- Return JSON responses directly except for very rare cases.

Routes answer: "Which handler should run for this HTTP request?"

## `internal/handlers`

Handlers are the HTTP layer. In other frameworks, these are often called controllers.

They are responsible for:

- Reading path parameters.
- Reading query parameters.
- Reading and validating JSON request bodies.
- Calling the next layer of the application.
- Choosing HTTP status codes.
- Returning JSON responses.

They should not:

- Contain large SQL queries.
- Own complex business rules.
- Know how tables are joined internally.
- Become large files that handle multiple unrelated features.

For simple endpoints, a handler may call a repository directly. For more complex workflows, a handler should call a service.

## `internal/models`

Models define the core data shapes used by the backend.

They are usually Go structs that match database records or important domain objects.

Example:

```go
type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
}
```

They are responsible for:

- Giving database/domain data a clear Go type.
- Defining shared field names and types.
- Providing tags for JSON responses and database mapping.

They should not:

- Connect to the database.
- Run queries.
- Handle HTTP requests.
- Contain route definitions.

Keep request/response-specific structs separate when the API shape differs from the database model. For example, a `CreateUserRequest` may include a password, but the `User` model should not expose one.

## `internal/repository`

Repositories contain database queries for a specific domain.

Example responsibilities:

- `CreateUser`
- `GetUserByEmail`
- `ListPlatforms`
- `AssignRoleToUser`

They are responsible for:

- Executing SQL queries.
- Scanning database rows into models.
- Returning models and errors to handlers or services.
- Keeping table-specific database details out of the HTTP layer.

They should not:

- Know about Gin.
- Read request bodies.
- Choose HTTP status codes.
- Return JSON responses.
- Contain broad business workflows that combine many unrelated actions.

Repositories answer: "How does this application read or write this data in the database?"

## `internal/middleware`

Middleware contains logic that runs before or around handlers.

Common examples:

- Authentication.
- Authorization.
- Request logging.
- CORS, if configured manually.
- Rate limiting.

Middleware is responsible for:

- Checking request-level concerns before the handler runs.
- Adding values to request context when needed.
- Stopping the request early when it is not allowed.

It should not:

- Contain feature-specific business logic.
- Run unrelated database workflows.
- Return normal feature responses.

## Future Packages

These packages are not required for every feature, but they become useful as the backend grows.

### `internal/services`

Services contain business workflows that combine multiple steps.

Example responsibilities:

- Registering a user.
- Hashing a password.
- Creating a platform and assigning the owner role.
- Checking whether a user can perform an action.

Services sit between handlers and repositories when the logic is more than a simple database call.

## Request Flow

A typical request should eventually flow like this:

```text
HTTP request
-> route
-> middleware
-> handler
-> service
-> repository
-> database
-> JSON response
```

For simple endpoints, the service layer can be skipped:

```text
HTTP request
-> route
-> handler
-> repository
-> database
-> JSON response
```

## Migrations

The `migrations` directory contains SQL files that define database schema changes.

Each migration should have:

- An `.up.sql` file for applying the change.
- A `.down.sql` file for rolling the change back.

Do not create tables manually in application startup code. Schema changes should be captured in migrations.

## Scripts

The `scripts` directory contains helper scripts for local development.

The migration script should use the same environment variable contract as the API config so local setup stays consistent.

## General Rules

- Keep `main.go` small.
- Keep route registration in `internal/routes`.
- Keep HTTP request/response logic in `internal/handlers`.
- Keep database connection setup in `internal/database`.
- Keep table/domain structs in `internal/models`.
- Add repositories when SQL queries start appearing.
- Add services when workflows involve business rules or multiple repositories.
- Do not let one package become a dumping ground.
