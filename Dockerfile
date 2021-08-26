FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod download

ENTRYPOINT go run main.go

EXPOSE 8080
