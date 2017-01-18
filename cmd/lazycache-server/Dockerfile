FROM amarburg/docker-golang-ffmpeg:wheezy-1.8

RUN go get github.com/amarburg/go-lazycache

RUN go install github.com/amarburg/go-lazycache/cmd/lazycache-server/

ENTRYPOINT /go/bin/lazycache-server

# Document that the service listens on port 5000
EXPOSE 5000
