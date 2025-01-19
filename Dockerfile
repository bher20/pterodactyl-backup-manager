FROM golang:1.20 AS build-stage

ENV \
    APP_DIR=/opt/app

WORKDIR ${APP_DIR}

RUN \
    mkdir -p ${APP_DIR}

COPY go.mod go.sum ${APP_DIR}/

RUN \
    go mod download

COPY ./internal ${APP_DIR}/internal
COPY ./cmd ${APP_DIR}/cmd
COPY ./*.go ${APP_DIR}/


RUN CGO_ENABLED=0 GOOS=linux go build -o /pterodactyl-backup-manager


FROM alpine:3 AS build-release-stage

COPY --from=build-stage /pterodactyl-backup-manager /pterodactyl-backup-manager

RUN apk add --no-cache --update curl ca-certificates openssl git tar bash sqlite fontconfig \
    && adduser --disabled-password --home /home/container container

USER container
ENV  USER=container HOME=/home/container

WORKDIR /home/container

COPY ./scripts/entrypoint.sh /entrypoint.sh
COPY ./metadata.json /metadata.json

CMD ["/bin/bash", "/entrypoint.sh"]

LABEL \
  MAINTAINER="Bherville, <support@bherville.com>" \
  org.opencontainers.image.source=https://github.com/bherville/pterodactyl-backup-manager
    