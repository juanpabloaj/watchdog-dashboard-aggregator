# Watchdog dashboard

Run the service

    go run ./cmd/server

Available endpoint

    GET /dashboard/{id} - Get user dashboard

Test manually

    curl http://localhost:8080/dashboard/6

Environment variables example are available in .env.example file.

Run tests

    go test -v ./...

Specific tests related to API timeout

* TestGetUserTimeoutWithSleep
* TestGetUserTimeoutWithMock

`TestGetUserTimeoutWithSleep` is skipped because it is slow, if you want to run it, comment the skip line.

Those tests are in the file internal/infrastructure/dummyjson/get_user_test.go.

## Notes

* User and todos are fetched concurrently using goroutines and sync.WaitGroup.
* I used the standard library because the service is small, the only external dependency is `go-cmp` to compare structs in tests.
* The containerized version logs the build date, the commit hash, and the env variables, it is useful for debugging and monitoring.

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
