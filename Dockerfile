FROM alpine:3.6
MAINTAINER Rohith Jayawardene <rohith.jayawardene@appvia.io>

RUN apk add ca-certificates --update --no-cache

ADD bin/controller /controller

ENTRYPOINT ["/controller"]
