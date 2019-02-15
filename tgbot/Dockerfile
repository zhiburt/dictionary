FROM golang

MAINTAINER zhiburt

ADD . /go/src/github.com/dictionary/tgbot

RUN go install github.com/dictionary/tgbot

ENTRYPOINT /go/bin/tgbot