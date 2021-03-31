FROM golang:1.14-alpine

ENV CGO_ENABLED=false

WORKDIR /app

COPY go.* ./

RUN go mod download && go get -u github.com/cosmtrek/air

COPY . .
