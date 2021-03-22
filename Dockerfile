FROM golang:1.16-alpine

WORKDIR /app

COPY . /app

RUN apk update \
    && apk --no-cache --update add build-base alpine-sdk

RUN go build -o main

RUN ls | sed 's/^main//' | xargs rm -rf

ENTRYPOINT ["./main"]
