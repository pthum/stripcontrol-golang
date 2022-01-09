#!/bin/sh
docker build --tag stripcontrol-go -f build/package/Dockerfile ./
# docker run --name stripcontrol-go --rm -p 8080:8080 stripcontrol-go