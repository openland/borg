FROM ubuntu:16.04

# Golang Paths
ENV GOVERSION 1.10
ENV GOROOT /opt/go
ENV GOPATH /root/.go

# Basic dependencies
RUN apt-get update && apt-get -y install software-properties-common python-software-properties gnupg-agent && apt-get update

# Additional repos
RUN add-apt-repository -y ppa:ubuntugis/ubuntugis-unstable

# Installing all dependencies
RUN apt-get -y install \
    wget \
    curl \
    git \
    cmake \
    gnupg-agent \
    gdal-bin \
    gcc \
    build-essential \
    libcgal-qt5-dev \
    libcgal-dev \
    libgdal-dev

# Installing Go
RUN cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

# Building prepair
RUN cd /opt && \
    git clone --depth=1 https://github.com/statecrafthq/prepair.git && \
    cd /opt/prepair && \
    cmake . && \
    make && \
    mv prepair /usr/bin/ && \
    rm -fr /opt/prepair

# Copying sources
RUN mkdir -p /root/.go/src/github.com/statecrafthq/borg/
COPY . /root/.go/src/github.com/statecrafthq/borg/

# Go Dependencies
RUN go get \
    gopkg.in/kyokomi/emoji.v1 \
    github.com/urfave/cli \
    github.com/twpayne/go-geom \
    github.com/buger/jsonparser \
    gopkg.in/cheggaaa/pb.v1 \
    cloud.google.com/go/storage

# Building Go
RUN cd /root/.go/src/github.com/statecrafthq/borg/ && go build && mv borg /usr/bin/

ENTRYPOINT ["/usr/bin/borg"]