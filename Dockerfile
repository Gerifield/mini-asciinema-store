FROM golang:1.14-alpine AS builder

RUN apk update && apk add ca-certificates

ADD . /build
WORKDIR /build

ENV CGO_ENABLED=0

RUN go test -v ./...
RUN go build src/cmd/server/server.go


FROM alpine:3.11

COPY --from=builder /build/server /usr/local/bin/
RUN mkdir /uploads

ENTRYPOINT ["/usr/local/bin/server"]

EXPOSE 8080