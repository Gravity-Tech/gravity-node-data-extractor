FROM golang:1.16-alpine

WORKDIR /app

COPY . /app

RUN apk add build-base
RUN go build -o main

RUN ls | sed 's/^main//' | xargs rm -rf

ENTRYPOINT ["./main"]
