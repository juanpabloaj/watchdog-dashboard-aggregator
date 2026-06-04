# Watchdog dashboard

Run the service

    go run ./cmd/server

Available endpoint

    GET /dashboard/{id} - Get user dashboard

Test manually

    curl http://localhost:8080/dashboard/6

Environment variable examples are available in `.env.example`.

Run tests

    go test -v ./...

Timeout related tests.

* TestGetUserTimeoutWithSleep
* TestGetUserTimeoutWithMock

`TestGetUserTimeoutWithSleep` is skipped because it is slow, to run it, remove the skip line.

Those tests are in the file internal/infrastructure/dummyjson/get_user_test.go.

## Design Notes

* I used a hexagonal structure with ports and adapters to reduce the coupling between transport, application logic, and external services.

```
    ├── application
    │   └── dashboard
    │       ├── dashboard.go
    ├── domain
    │   ├── dashboard.go
    │   ├── todo.go
    │   └── user.go
    ├── infrastructure
    │   └── dummyjson
    │       ├── dummyjson.go
    └── interfaces
        └── httphandler
            ├── dashboardhandler.go
```

* The HTTP Client has a timeout of 2 seconds by default, it can be changed by an environment variable.
* User and todos are fetched concurrently using goroutines and `sync.WaitGroup`.
* I used the standard library because the service is small. The only external dependency is `go-cmp` to compare structs in tests.
* All the packages have tests. I focused on happy paths and requirement cases, I didn't cover all the scenarios (I didn't cover all the HTTP scenarios, etc).
* The containerized version logs the build date, the commit hash, and the env variables, it is useful for debugging and monitoring.

## Requirement coverage

* User Details: Combine firstName and lastName.
* Status Calculation:
    * If the user is age > 50, their status is "Veteran".
    * Otherwise, their status is "Rookie".
* Todo Summary:
    * Filter the todos to find only those where completed: false.
    * Return the count of pending tasks and the title of the first pending task
(if any).

    internal/application/dashboard/aggregate_test.go

* The external API is "slow" (simulate this or assume it). Your client must have a strict timeout of 2 seconds.

    internal/infrastructure/dummyjson/get_user_test.go

* If fetching Todos fails but User succeeds, return the User data with a "Todos Unavailable" warning (Partial Failure handling).

    internal/application/dashboard/get_dashboard_test.go

* Expected JSON response.

    internal/interfaces/httphandler/get_dashboard_test.go

## Docker and docker-compose

    docker-compose up -d --build

    docker-compose logs --follow

    docker-compose down


## Out of scope

* Authentication
* CORS
* Metrics and distributed tracing
* Graceful shutdown
* Kubernetes endpoints, health checks (liveness and readiness)
* Kubernetes files
* CI/CD files
