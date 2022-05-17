FROM golang:alpine as build-env

RUN apk update
RUN apk add --no-cache git

ADD . /app/
WORKDIR /app

ENV GOPROXY https://goproxy.io
RUN go mod download
RUN go build -o main .

FROM alpine:latest
COPY --from=build-env /app/main /app/main

RUN apk update
RUN apk add --no-cache tzdata

EXPOSE 9000

ENTRYPOINT ["/app/main"]