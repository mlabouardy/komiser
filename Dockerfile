FROM alpine:3.9.4
MAINTAINER mlabouardy <mohamed@labouardy.com>

ENV VERSION 2.4.0
ENV PORT 3000
ENV DURATION 30

RUN apk update && apk add curl
RUN curl -L https://s3.us-east-1.amazonaws.com/komiser/$VERSION/linux/komiser -o /usr/bin/komiser && \
    chmod +x /usr/bin/komiser && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    mkdir ~/.aws && touch ~/.aws/credentials

EXPOSE $PORT
ENTRYPOINT ["komiser", "start"]
