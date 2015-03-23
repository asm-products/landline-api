FROM golang:1.3-onbuild

ADD  . /go/src/github.com/asm-products/landline-api
WORKDIR /go/src/github.com/asm-products/landline-api

ENV GOPATH /go/src/github.com/asm-products/landline-api/Godeps/_workspace:$GOPATH
RUN  go install -v -a

ENTRYPOINT [ "./landline-api" ]
