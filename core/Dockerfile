FROM golang:1.18.9-alpine3.17 AS build-env

ENV GO111MODULE=on

WORKDIR /workspace/core

COPY ./manager/ ./manager
COPY ./core.go ./core.go
COPY ./go.mod ./go.mod

RUN go mod tidy -v

RUN go build -o ./application .

FROM alpine:3.17

WORKDIR /app

COPY --from=build-env /workspace/core/application /core

ENV MANAGER_PORT=8080

ENV NATS_URL=localhost
ENV NATS_USERNAME=""
ENV NATS_PASSWORD=""
EXPOSE ${MANAGER_PORT}
EXPOSE ${PROXY_PORT}
EXPOSE ${GATEWAY_PORT}

CMD /core \
    --manager_port=${MANAGER_PORT} \
    --nats_url=${NATS_URL} --nats_username=${NATS_USERNAME} --nats_password=${NATS_PASSWORD}