FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOVERSION 1.10
ENV GOROOT /opt/go
ENV GOPATH /root/.go

# Installing Go
RUN apt-get update && apt-get -y install \
    wget \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*
RUN cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

# Import sources
RUN mkdir -p /root/.go/src/github.com/statecrafthq/borg/
COPY . /root/.go/src/github.com/statecrafthq/borg/

# Go Dependencies
RUN go get \
    gopkg.in/kyokomi/emoji.v1 \
    github.com/urfave/cli \
    github.com/twpayne/go-geom \
    github.com/buger/jsonparser \
    gopkg.in/cheggaaa/pb.v1
RUN cd /root/.go/src/github.com/statecrafthq/borg/ && go install

ENTRYPOINT ["/root/.go/bin/borg"]