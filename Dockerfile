FROM golang:1.12-alpine as builder

RUN apk add git

WORKDIR /code

ADD . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOCACHE=/tmp/.cache/go-build
ENV GOPATH=/tmp/go

RUN go build -o goOrderAPI

FROM alpine:3.10.1

RUN mkdir /code
RUN apk add --no-cache \
        libc6-compat \
        ca-certificates

WORKDIR /code

COPY --from=builder /code/goOrderAPI .

EXPOSE 8000
CMD ./goOrderAPI