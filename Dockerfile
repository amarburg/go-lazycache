FROM amarburg/golang-ffmpeg:wheezy-1.8

#
## Build the local repo in Dockerfile
#
ADD . $GOPATH/src/github.com/amarburg/go-lazycache

RUN go get -v github.com/amarburg/go-lazycache/app
RUN go build -v github.com/amarburg/go-lazycache/app
RUN go install -v github.com/amarburg/go-lazycache/app

VOLUME ["/srv/image_store"]

ENV LAZYCACHE_PORT=8080

## Strangely, it's installing the binary to $GOPATH
CMD $GOPATH/lazycache
EXPOSE 8080
