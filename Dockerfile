FROM golang:1.3

COPY . /go/src/github.com/crosbymichael/dockersql
WORKDIR /go/src/github.com/crosbymichael/dockersql

RUN go get -d && go build .

ENTRYPOINT ["./dockersql"]
