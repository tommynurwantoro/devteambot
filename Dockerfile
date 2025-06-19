FROM golang:1.23-alpine AS builder

LABEL maintainer="Tommy Nurwantoro <tommy.nurwantoro@gmail.com>"

WORKDIR /go/src/app
ADD . /go/src/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main


FROM alpine:latest

LABEL maintainer="Tommy Nurwantoro <tommy.nurwantoro@gmail.com>"

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=builder /go/src/app/config.yaml /app/config.yaml
COPY --from=builder /go/src/app/main /app/main

EXPOSE 9050
ENTRYPOINT [ "./main", "service" ]