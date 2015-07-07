FROM golang
MAINTAINER Douézan-Grard Guillaume - Quorums

RUN go get github.com/quorumsco/users

ADD . /go/src/github.com/quorumsco/users

WORKDIR /go/src/github.com/quorumsco/users

RUN \
  go get -u && \
  go build

EXPOSE 8080

ENTRYPOINT ["./users"]
