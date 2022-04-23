#!/usr/bin/env bash
# docker run --name authz -p 192.168.1.1:3001:3001 -p -d "$REGISTRY_PREFIX"/"$PROJECT_PREFIX"/authz-"$ARCH":"$VERSION"
docker run --name authz -p 192.168.1.1:3001:3001 --net basic-service --ip 192.168.10.3 -d 192.168.1.200:80/basic-service/authz-amd64:0.0.1