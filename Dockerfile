FROM golang:1.4

RUN apt-get update
RUN apt-get install -y inotify-tools psmisc && \
    wget https://raw.github.com/alexedwards/go-reload/master/go-reload && \
        chmod +x go-reload && \
        mv go-reload /usr/local/bin/

COPY  . /go/src/github.com/asm-products/landline-api
WORKDIR /go/src/github.com/asm-products/landline-api

ENV GOPATH /go/src/github.com/asm-products/landline-api/Godeps/_workspace:$GOPATH
ENV PATH /go/src/github.com/asm-products/landline-api/Godeps/_workspace/bin:$PATH
RUN  go install -v -a

EXPOSE 3000

CMD [ "go-reload", "main.go" ]
