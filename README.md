# Watchdog dashboard

Run service

    go run ./cmd/server

Available endpoint

    GET /dashboard/{1} - Get user dashboard

Test manually

    curl 0.0.0.0:8080/dashboard/1

Environment variables example available in .env.example file.

Run tests

    go test -v ./...

Specific tests related to API timeout

* TestGetUserTimeoutWithSleep
* TestGetUserTimeoutWithMock

`TestGetUserTimeoutWithSleep` is skipped because it is slow, if you want to run it, comment the skip line.

Those tests are in the file internal/infrastructure/dummyjson/get_user_test.go.

## Docker and docker-compose

    docker-compose up -d --build

    docker-compose logs --follow

    docker-compose down


## Out of scope

* Authentication
* CORS
* Metrics and Distributed tracing
* Graceful shutdown
* Kubernetes endpoints, health checks (liveness and readiness)
* Kubernetes files
* CI/CD files
