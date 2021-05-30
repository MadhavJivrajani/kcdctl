FROM golang:alpine 

MAINTAINER Madhav Jivrajani <madhav.jiv@gmail.com>

ADD . /src
WORKDIR /src
RUN go build -o lb loadbalancer.go

ENV containerPort=8080
ENV targetPort=8080
EXPOSE $containerPort

ENV serviceSelector=""
ENV serviceSelectorValue=""

ENTRYPOINT ["./lb"]
