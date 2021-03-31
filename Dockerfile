FROM golang:1.14-alpine

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.* ./

RUN go mod download && go get -u github.com/cosmtrek/air

COPY . .
