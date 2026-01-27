FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

RUN GOBIN=/app go install github.com/rubenv/sql-migrate/...@v1.6.1

FROM alpine:3.23

WORKDIR /app/

COPY --from=builder /app/sql-migrate /app/sql-migrate
COPY script/db/migrations /app/script/db/migrations
#COPY config/dbconfig.yml /app/config/dbconfig.yml

ENTRYPOINT ["./sql-migrate", "up", "-config", "./config/dbconfig.yml"]
