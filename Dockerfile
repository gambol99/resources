FROM alpine:3.7
MAINTAINER Rohith Jayawardene <rohith.jayawardene@appvia.io>

RUN apk add ca-certificates curl --no-cache && \
    adduser -D controller

ADD bin/controller /controller

USER 1000

ENTRYPOINT ["/controller"]
