FROM golang

RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template-consumer
RUN mkdir -p /go/src/github.com/wednesday-solutions/libs
RUN mkdir -p /go/src/github.com/wednesday-solutions/utils

RUN [ "go", "install", "github.com/rubenv/sql-migrate/...@latest"]

ADD ./libs /go/src/github.com/wednesday-solutions/libs
ADD ./utils /go/src/github.com/wednesday-solutions/utils
ADD ./consumer /go/src/github.com/wednesday-solutions/go-template-consumer

WORKDIR /go/src/github.com/wednesday-solutions/libs
RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
WORKDIR /go/src/github.com/wednesday-solutions/utils
RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
WORKDIR ../go-template-consumer

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
RUN go build -o ./ ./cmd/server/main.go
CMD ["bash", "./migrate-and-run.sh"]
EXPOSE 9000