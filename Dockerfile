FROM golang:1.16-alpine

WORKDIR /app

COPY . /app

RUN go build -o main

RUN ls | sed 's/^main//' | xargs rm -rf

ENTRYPOINT ["./main"]
