FROM amarburg/golang-ffmpeg:wheezy-1.8

RUN go get github.com/amarburg/go-lazycache-app
RUN go install github.com/amarburg/go-lazycache-app

CMD ["/go/bin/go-lazycache-app"]
EXPOSE 8080

ENV LAZYCACHE_PORT=8080
