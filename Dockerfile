FROM golang
MAINTAINER Douézan-Grard Guillaume - Quorums

ADD . /go/src/github.com/quorumsco/users

WORKDIR /go/src/github.com/quorumsco/users

RUN \
  go get && \
  go build

EXPOSE 8080

ENTRYPOINT ["./users"]
