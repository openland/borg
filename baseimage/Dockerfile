FROM ubuntu:16.04

# Golang Paths
ENV GOVERSION 1.10
ENV GOROOT /opt/go
ENV GOPATH /root/.go

# Basic dependencies
RUN apt-get update && \
    apt-get -y install software-properties-common python-software-properties gnupg-agent && \
    apt-get update && \
    add-apt-repository -y ppa:ubuntugis/ubuntugis-unstable && \
    apt-get -y install \
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
    libgdal-dev && \
    cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH && \
    cd /opt && \
    git clone --depth=1 https://github.com/statecrafthq/prepair.git && \
    cd /opt/prepair && \
    cmake . && \
    make && \
    mv prepair /usr/bin/ && \
    rm -fr /opt/prepair && \
    rm -rf /var/lib/apt/lists/*