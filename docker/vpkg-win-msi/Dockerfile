FROM ubuntu:16.04

MAINTAINER Craig Smith <craig.smith@768bit.com>

RUN apt-get update && apt-get upgrade -y
RUN curl -sL https://deb.nodesource.com/setup_10.x -o nodesource_setup.sh
RUN apt-get install msitools nodejs

RUN npm install -g @768bit/msi-packager

ENTRYPOINT ["msi-packager"]
