# syntax=docker/dockerfile:1

# This is the base image for the different services

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd
COPY config ./config
