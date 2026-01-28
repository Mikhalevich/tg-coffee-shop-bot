FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/manager cmd/manager/main.go

FROM alpine:3.23

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/manager /app/manager
#COPY --from=builder /app/config/config-http-manager.yaml /app/config/config-http-manager.yaml

ENTRYPOINT ["./manager", "-config", "./config/config-http-manager.yaml"]
