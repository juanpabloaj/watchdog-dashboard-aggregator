FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags "-X 'main.gitHash=$(git describe --tags --always --dirty)' -X 'main.buildDate=$(date -Isecond)'" \
  -o application ./cmd/server/*.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/application .

EXPOSE 8080

CMD ["/app/application"]
