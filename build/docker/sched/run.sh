#!/usr/bin/env bash
# docker run --name sched -p 192.168.1.1:3001:3001 -p -d "$REGISTRY_PREFIX"/"$PROJECT_PREFIX"/sched-"$ARCH":"$VERSION"
docker run --name sched -p 192.168.1.1:3001:3001 --net basic-service --ip 192.168.10.6 -d 192.168.1.200:80/basic-service/sched-amd64:0.0.1