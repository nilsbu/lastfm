# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd

RUN go build -o /lastfm-wrkr ./cmd/lastfm-wrkr

CMD /lastfm-wrkr /data
