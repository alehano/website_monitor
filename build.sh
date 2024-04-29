#!/bin/sh


# docker buildx create --name mybuilder --use
docker buildx use mybuilder
docker buildx build --platform linux/amd64,linux/arm64 -t akhalyapin/website_monitor:latest --push .
