FROM golang

RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template-producer
RUN mkdir -p /go/src/github.com/wednesday-solutions/libs
RUN mkdir -p /go/src/github.com/wednesday-solutions/utils

RUN [ "go", "install", "github.com/rubenv/sql-migrate/...@latest"]
RUN [ "go", "install", "github.com/volatiletech/sqlboiler@v3.7.1"]
RUN [ "go", "install", "github.com/99designs/gqlgen@v0.14.0"]

ADD ./libs /go/src/github.com/wednesday-solutions/libs
ADD ./utils /go/src/github.com/wednesday-solutions/utils
ADD ./producer /go/src/github.com/wednesday-solutions/go-template-producer

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
WORKDIR ../go-template-producer

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
RUN go build -o ./ ./cmd/server/main.go
CMD ["bash", "./migrate-and-run.sh"]
EXPOSE 9001
