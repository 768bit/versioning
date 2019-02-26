#!/usr/bin/env bash


IMAGE_REGISTRY_PATH="registry.gitlab.768bit.com/github/vpkg/docker"

##docker build -t $IMAGE_REGISTRY_PATH/vpkg-macos docker/vpkg-macos

docker build -t $IMAGE_REGISTRY_PATH/vpkg-win-msi docker/vpkg-win-msi
