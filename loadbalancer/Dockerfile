# FROM golang:alpine 
FROM scratch

MAINTAINER Madhav Jivrajani <madhav.jiv@gmail.com>

# ADD . /src
# WORKDIR /src
# RUN go build -o lb loadbalancer.go
ADD lb lb
ENV exposedPort=9000
ENV targetPort=8080
EXPOSE $exposedPort

ENV serviceSelector=""
ENV serviceSelectorValue=""

ENTRYPOINT ["/lb"]
