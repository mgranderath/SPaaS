FROM gliderlabs/alpine:3.4

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app
ADD release/SPaaS_server /app/

EXPOSE 8080

ENTRYPOINT ["./SPaaS_server"]