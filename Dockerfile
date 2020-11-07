FROM golang:1.15.4-alpine3.12 AS builder

RUN echo $'\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/main\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/testing\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

COPY . /build
WORKDIR /build
RUN go build -o /usr/bin/E

FROM alpine:3.12

RUN echo $'\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/main\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/testing\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

ENV UID=1255 \
    GID=1255

RUN apk add --no-cache su-exec ca-certificates olm bash curl

COPY --from=builder /usr/bin/E /usr/bin/E
COPY --from=builder /build/docker-run.sh /docker-run.sh

VOLUME /data

CMD ["/docker-run.sh"]
