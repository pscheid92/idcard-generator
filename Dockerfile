FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY ./idcard-generator /idcard-generator
COPY templates /templates

EXPOSE 8080
ENTRYPOINT ["/idcard-generator"]
