FROM debian
MAINTAINER Maksym Prokpov <mprokopov@gmail.com>

RUN apt-get update && apt-get install -y git-core

WORKDIR /app

COPY getversion.linux /usr/sbin/getversion

CMD ["getversion"]
