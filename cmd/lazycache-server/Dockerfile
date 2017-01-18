FROM amarburg/docker-golang-ffmpeg:wheezy-1.8

RUN go get github.com/amarburg/go-lazycache

RUN go install github.com/amarburg/go-lazycache/cmd/lazycache-server/

COPY service_account.json /go/
ENV GOOGLE_APPLICATION_CREDENTIALS=/go/service_account.json

ENTRYPOINT /go/bin/lazycache-server

# Document that the service listens on port 5000
EXPOSE 5000
