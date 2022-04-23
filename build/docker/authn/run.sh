#!/usr/bin/env bash
# docker run --name authn -p 192.168.1.1:3001:3001 -p -d "$REGISTRY_PREFIX"/"$PROJECT_PREFIX"/authn-"$ARCH":"$VERSION"
docker run --name authn -p 192.168.1.1:3001:3001 --net basic-service --ip 192.168.10.2 -d 192.168.1.200:80/basic-service/authn-amd64:0.0.1