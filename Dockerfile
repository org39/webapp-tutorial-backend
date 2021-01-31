# builder container
FROM golang:1.15-alpine3.12 AS builder
RUN apk update && apk add --no-cache git make

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build/server

#####################
FROM alpine:3.12

# install SSL root certificates
RUN apk update && apk add ca-certificates && \
	rm -rf /var/cache/apk/* && \
	update-ca-certificates

RUN adduser -S -H -u 3939 u u
USER 3939
WORKDIR /opt
COPY --from=builder /build/server /opt

CMD ["/opt/server"]
