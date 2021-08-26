FROM golang:latest

RUN apt update
RUN apt install uuid-runtime
RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod download

ENTRYPOINT go run main.go

EXPOSE 8080
