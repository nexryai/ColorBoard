FROM alpine:edge
WORKDIR /workspace

RUN apk add --no-cache ca-certificates go go-task go-task-task cargo wasm-pack alpine-sdk g++ build-base cmake clang libressl-dev python3 nodejs-current yarn \
 && addgroup -g 630 app \
 && adduser -u 630 -G app -D -h /app app