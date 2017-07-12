FROM golang

ADD . /go/src/github.com/baskeboler/registry
WORKDIR /go/src/github.com/baskeboler/registry

RUN go get && go install

ENTRYPOINT /go/bin/registry

ENV PORT=8080

EXPOSE 8080
