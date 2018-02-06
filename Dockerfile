FROM golang:1.9.2

# nodejs + apidocjs
RUN apt update && \
    curl --silent --location https://deb.nodesource.com/setup_6.x | bash - && \
    apt install -y nodejs && \
    npm install apidoc -g

ENV APP_SOURCE $GOPATH/src/github.com/av-ast/sms_sender

RUN mkdir -p $APP_SOURCE
ADD . $APP_SOURCE
WORKDIR $APP_SOURCE

RUN make

EXPOSE 8000

CMD ["bash", "-c", "sms_sender"]
