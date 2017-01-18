FROM amarburg/docker-golang-ffmpeg:wheezy-1.8

RUN go get github.com/amarburg/go-lazycache

RUN go install github.com/amarburg/go-lazycache/cmd/server/lazycache.go

ENTRYPOINT /go/bin/lazycache

# Document that the service listens on port 8080.
EXPOSE 80
