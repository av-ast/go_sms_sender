FROM golang:1.9.2

ENV APP_SOURCE $GOPATH/src/github.com/av-ast/sms_sender

RUN mkdir -p $APP_SOURCE
ADD . $APP_SOURCE
WORKDIR $APP_SOURCE

RUN make

EXPOSE 8000

CMD ["bash", "-c", "sms_sender"]
