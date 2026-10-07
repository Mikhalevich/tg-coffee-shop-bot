FROM alpine:3.24

RUN apk --no-cache add curl

ARG CONNECT_HOST=localhost:8083
ENV CONNECT_HOST_ENV=${CONNECT_HOST}

WORKDIR /app

COPY script/docker/debezium-config.json /app/debezium-config.json

CMD curl -s -S -XPOST -H Accept:application/json -H Content-Type:application/json http://${CONNECT_HOST_ENV}/connectors/ -d @/app/debezium-config.json
