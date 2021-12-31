FROM golang:1.17 AS build

WORKDIR /hongbao-ms

COPY build/Release.key /Release.key
RUN apt-get update && \
        apt-get -y --no-install-recommends install gnupg2 && \
        echo "deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/ ./" >> /etc/apt/sources.list && \
        apt-key add /Release.key && \
        apt-get -y --no-install-recommends install libczmq-dev

COPY ./ /hongbao-ms
RUN GOPROXY=https://goproxy.io,direct go build -o /usr/local/bin/msd cmd/msd/msd.go

CMD ["msd"]