FROM ubuntu:focal

RUN apt-get update && apt-get install -y curl
RUN cd /usr/local/bin && curl -o kuma-0.7.1-ubuntu-amd64.tar.gz -sSL https://kong.bintray.com/kuma/kuma-0.7.1-ubuntu-amd64.tar.gz && tar xvfz kuma-0.7.1-ubuntu-amd64.tar.gz && rm -f kuma-0.7.1-ubuntu-amd64.tar.gz && mv -v ./kuma-0.7.1/bin/* . && rm -rf ./kuma-0.7.1
