FROM golang:1.14-alpine

WORKDIR /app

COPY . /app

RUN export GOOS=linux && export GOARCH=amd64 \
    && ./buildhelper.sh && go build -o main main.go

ENTRYPOINT ["./main"]