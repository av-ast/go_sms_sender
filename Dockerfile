FROM golang:1.9.2

ENV APP_SOURCE $GOPATH/src/github.com/av-ast/sms_sender

RUN mkdir -p $APP_SOURCE
ADD . $APP_SOURCE
WORKDIR $APP_SOURCE

RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    go install

EXPOSE 8000

CMD ["bash", "sms_sender"]
