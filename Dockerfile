# === BUILD STAGE === #
FROM golang:1.12-alpine as build

ARG ACCESS_TOKEN

RUN apk add --no-cache git

WORKDIR /srv/app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test -v ./...
RUN go build -ldflags="-w -s" -o build

# === RUN STAGE === #
FROM alpine as run

RUN apk update \
        && apk upgrade \
        && apk add --no-cache ca-certificates \
        && update-ca-certificates \
        && rm -rf /var/cache/apk/*
        
WORKDIR /srv/app
COPY --from=build /srv/app/build /srv/app/build

ENV LISTEN_ADDRESS 0.0.0.0:8080
ENV GIN_MODE release

CMD ["/srv/app/build"]
