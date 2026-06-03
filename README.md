# Watchdog dashboard

Run test

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

* CORS
* Authentication
* Metrics and Distributed tracing
