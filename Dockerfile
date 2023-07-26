FROM golang:1.15.0-alpine AS builder

LABEL maintainer="Tommy Nurwantoro <tommy.nurwantoro@gmail.com>"

# Install shell requirement
RUN apk update && apk upgrade && apk add --no-cache git tzdata ca-certificates

ADD . /app/devteambot
WORKDIR /app/devteambot

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o devteambot main.go

FROM alpine:latest
WORKDIR /app/devteambot

COPY --from=builder /app/devteambot/devteambot /app/devteambot/devteambot
COPY --from=builder /app/devteambot/config.yaml /app/devteambot/config.yaml

EXPOSE 8000
ENTRYPOINT [ "./devteambot", "svc", "--config=config.yaml" ]