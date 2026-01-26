FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/msgconsumer cmd/msgconsumer/main.go

FROM alpine:3.23

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/msgconsumer /app/msgconsumer
#COPY --from=builder /app/config/config-msgconsumer.yaml /app/config-msgconsumer.yaml

ENTRYPOINT ["./msgconsumer", "-config", "config-msgconsumer.yaml"]
