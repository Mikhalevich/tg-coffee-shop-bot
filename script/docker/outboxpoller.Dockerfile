FROM golang:1.26-alpine3.24 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/outboxpoller cmd/outboxpoller/main.go

FROM alpine:3.24

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/outboxpoller /app/outboxpoller
#COPY --from=builder /app/config/config-outbox-poller.yaml /app/config/config-outbox-poller.yaml

ENTRYPOINT ["./outboxpoller", "-config", "./config/config-outbox-poller.yaml"]
