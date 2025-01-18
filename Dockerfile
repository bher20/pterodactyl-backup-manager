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



FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /pterodactyl-backup-manager /pterodactyl-backup-manager

USER nonroot:nonroot

ENTRYPOINT [ "/pterodactyl-backup-manager", "serve" ]