FROM golang:alpine

MAINTAINER Madhav Jivrajani <madhav.jiv@gmail.com>

ADD . /src
WORKDIR /src
RUN go build -o server server.go

EXPOSE 8080

ENTRYPOINT ["./server"]
