# Build the binary with minimal go image
FROM golang:1.26.4-alpine AS builder

# Select/assign working directory
WORKDIR /app

# Copy dependencies and download
COPY go.mod go.sum ./
RUN go mod download

# copy source code and build it
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY migrations/ ./migrations/
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

RUN apk --no-cache add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

EXPOSE 8000
CMD [ "./entrypoint.sh" ]
