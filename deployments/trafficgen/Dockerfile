# syntax=docker/dockerfile:1

FROM alpine

RUN apk add --no-cache bash curl

COPY run.sh /
RUN chmod +x /run.sh

ENTRYPOINT "/run.sh" "http://petstore:8080/api"