FROM golang:1.14 as build
WORKDIR /go/src/github.com/st-vasyl/chobot/
COPY . .
RUN go get -d -v github.com/BurntSushi/toml && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.9
LABEL MAINTAINER "Vasyl Stetsuryn <vasyl@vasyl.org"

ARG APK_FLAGS_COMMON="-q"
ARG APK_FLAGS_PERSISTANT="${APK_FLAGS_COMMON} --clean-protected --no-cache"

ENV LANG C.UTF-8
ENV TERM=xterm
USER root

RUN apk update && \
    apk add ${APK_FLAGS_PERSISTANT} \
            less \
            bash && \
    addgroup chobot && \
    adduser -u 1000 \
            -S \
            -D -G chobot \
            -h /home/chobot \
            -s /bin/bash \
            chobot && \
    mkdir -p /opt/chobot && \
    chown -R chobot:chobot /opt/chobot

COPY --from=build /go/src/github.com/st-vasyl/chobot/chobot /opt/chobot/chobot
USER chobot
WORKDIR /opt/chobot


CMD ["/opt/chobot/chobot"]
