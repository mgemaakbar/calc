FROM golang:1.8

RUN mkdir -p /go/src/calc
WORKDIR /go/src/calc

ADD . /go/src/calc

RUN go get -v