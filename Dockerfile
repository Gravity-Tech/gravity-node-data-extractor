FROM golang:1.15-buster

WORKDIR /app

COPY . /app

RUN go build -o main

RUN ls | sed 's/^main//' | xargs rm -rf

ENTRYPOINT ["./main"]
