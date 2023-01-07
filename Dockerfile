# syntax=docker/dockerfile:1

FROM golang:1.19 AS builder

WORKDIR /workspace

COPY . ./
RUN go mod download

ADD . .
RUN go build -o /ornn

FROM alpine
WORKDIR /
RUN apk add libgcc
COPY --from=builder /workspace/* /usr/local/bin
CMD ["ornn"]
