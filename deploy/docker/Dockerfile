FROM amarburg/golang-ffmpeg:wheezy-1.8

RUN go get -v github.com/amarburg/go-lazycache

## Hot patch local files into the repo
ADD *.go $GOPATH/src/amarburg/go-lazycache/app/

WORKDIR $GOPATH/src/github.com/amarburg/go-lazycache/app
RUN go get -v .
RUN go build -o lazycache .
RUN cp lazycache $GOPATH/

VOLUME ["/srv/image_store"]

ENV LAZYCACHE_PORT=8080

## Strangely, it's installing the binary to $GOPATH
CMD $GOPATH/lazycache
EXPOSE 8080
