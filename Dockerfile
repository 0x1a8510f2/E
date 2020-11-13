FROM golang:1.15.4-alpine3.12 AS builder

RUN echo $'\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/main\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/testing\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

COPY . /build
WORKDIR /build
RUN ./build.sh

FROM alpine:3.12

RUN echo $'\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/main\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/testing\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

ENV UID=1255 \
    GID=1255

RUN apk add --no-cache su-exec ca-certificates olm curl libcap

COPY --from=builder /build/E /usr/bin/E
COPY --from=builder /build/example-config.yaml /opt/E/example-config.yaml
COPY --from=builder /build/docker-run.sh /docker-run.sh

RUN setcap CAP_NET_BIND_SERVICE=+eip /usr/bin/E

VOLUME /data

CMD ["/docker-run.sh"]
