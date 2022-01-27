FROM golang

RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template-consumer

RUN go get -v github.com/rubenv/sql-migrate/... \
  github.com/volatiletech/sqlboiler \
  github.com/99designs/gqlgen

ADD . /go/src/github.com/wednesday-solutions/go-template-consumer
WORKDIR /go/src/github.com/wednesday-solutions/go-template-consumer

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
RUN go build -o ./ ./cmd/server/main.go
CMD ["bash", "./migrate-and-run.sh"]
EXPOSE 9000