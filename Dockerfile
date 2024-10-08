FROM alpine:edge as builder
WORKDIR /build

COPY . ./
RUN apk add --no-cache ca-certificates go go-task go-task-task cargo wasm-pack alpine-sdk g++ build-base cmake clang libressl-dev python3 nodejs-current yarn

RUN task build:client
RUN task build:server

FROM alpine:edge

RUN apk add --no-cache ca-certificates tini \
	&& addgroup -g 630 app \
	&& adduser -u 630 -G app -D -h /app app

USER app
WORKDIR /app

COPY --from=builder /build/cb-server ./
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/cb-server"]
