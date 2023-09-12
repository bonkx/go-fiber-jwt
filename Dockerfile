FROM golang:1.20.8-alpine3.18 as dev

RUN set -x -o pipefail \
    && apk update \
    && apk upgrade \
    && apk add build-base gcc g++ make git vips-dev vips-poppler pkgconf \
    && rm -rf /var/cache/apk/*

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

# manage dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .