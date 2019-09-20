FROM gliderlabs/alpine:3.9

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app
ADD release/SPaaS_server .
ADD release/post-receive ./util
ADD frontend/dist/. ./static

EXPOSE 8080

ENTRYPOINT ["./SPaaS_server"]