# syntax=docker/dockerfile:1

FROM golang:1.19 AS builder
RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

ADD . /go/src/github.com/gokch/ornn
WORKDIR /go/src/github.com/gokch/ornn
RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -o app .

FROM multiarch/ubuntu-core:arm64-bionic
WORKDIR /root/
COPY --from=builder /go/src/github.com/gokch/ornn .
CMD ["./ornn"]