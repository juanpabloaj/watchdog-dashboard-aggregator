FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags "-X 'main.gitHash=$(git describe --tags --always --dirty)' -X 'main.buildDate=$(date -Isecond)'" \
  -o application ./cmd/server/

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/application .

EXPOSE 8080

USER app

CMD ["/app/application"]
