FROM golang:alpine

ENV APP_PATH=$GOPATH/src/app

RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR $APP_PATH

COPY go.mod  go.sum ./

RUN go mod download

ADD . ./


EXPOSE 9001

RUN go build -o /usr/local/bin/consumer ./cmd/consumer/main.go && chmod +x /usr/local/bin/consumer

CMD ["/usr/local/bin/dockerize","-wait=tcp://rabbitmq:5672","-wait=tcp://kafka:9092","-wait=tcp://dispatcher:9090","-wait=tcp://persister:9090","-timeout=30s","-wait-retry-interval=5s","/usr/local/bin/consumer"]




