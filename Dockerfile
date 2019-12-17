FROM golang:1.11-alpine as builder
RUN apk --no-cache add curl ca-certificates git
WORKDIR /go/src/github.com/hr1sh1kesh/db-conn-sample
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY main.go  .
COPY Gopkg.toml .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o db-connect.o .

FROM alpine:latest
RUN apk --no-cache add ca-certificates git
WORKDIR /root/
COPY --from=builder /go/src/github.com/hr1sh1kesh/db-conn-sample/db-connect.o .
CMD ["./db-connect.o"]