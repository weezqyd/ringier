FROM golang:alpine

ENV APP_PATH=$GOPATH/src/app

WORKDIR $APP_PATH

COPY go.mod  go.sum ./

RUN go mod download

ADD . ./

EXPOSE 9090

RUN go build -o /usr/local/bin/dispatcher ./cmd/dispatcher/main.go && chmod +x /usr/local/bin/dispatcher

CMD ["/usr/local/bin/dispatcher"]




