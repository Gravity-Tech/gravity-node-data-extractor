FROM golang:1.16-alpine

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=0

RUN apk update \
	&& apk --no-cache --update add build-base

RUN go build -o main

RUN ls | sed 's/^main//' | xargs rm -rf

VOLUME /etc/extractor

ENTRYPOINT ["./main"]
