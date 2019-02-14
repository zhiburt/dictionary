FROM golang:latest

MAINTAINER zhiburt

ADD . /go/src/github.com/dictionary

RUN go install github.com/dictionary

ENTRYPOINT /go/bin/dictionary

EXPOSE 8081