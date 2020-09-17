FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY ./idcardgenerator /idcardgenerator
COPY templates /templates

EXPOSE 8080
ENTRYPOINT ["/idcardgenerator"]
