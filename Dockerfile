FROM amarburg/golang-ffmpeg:wheezy-1.8

#
## Build the local repo in Dockerfile
#
ADD . $GOPATH/src/github.com/amarburg/go-lazycache

WORKDIR $GOPATH/src/github.com/amarburg/go-lazycache/app
RUN go get -v .
RUN go build -o lazycache .
RUN cp lazycache $GOPATH/

VOLUME ["/srv/image_store"]

ENV LAZYCACHE_PORT=8080

## Strangely, it's installing the binary to $GOPATH
CMD $GOPATH/lazycache
EXPOSE 8080
