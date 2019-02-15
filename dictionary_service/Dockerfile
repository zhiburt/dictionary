FROM golang:latest

MAINTAINER zhiburt

ADD . /go/src/github.com/dictionary/dictionary_service

RUN go install github.com/dictionary/dictionary_service

ENTRYPOINT /go/bin/dictionary_service

#EXPOSE 8081