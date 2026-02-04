# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder
RUN apk add --no-cache make

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY Makefile ./
RUN make deps

COPY ./ ./
RUN make build

FROM scratch

COPY --from=builder /app/server /server

EXPOSE 8080

CMD [ "/server" ]