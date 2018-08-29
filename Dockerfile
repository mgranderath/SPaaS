FROM gliderlabs/alpine:3.4

WORKDIR /app
ADD release/SPaaS_server /app/

EXPOSE 8080

ENTRYPOINT ["./SPaaS_server"]