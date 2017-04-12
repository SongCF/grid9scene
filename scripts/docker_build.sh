#!/bin/bash

git checkout .
make clean

docker login -usoa -pAabb0011 139.198.2.55
docker pull 139.198.2.55/soalib/golang:1.8
docker run --rm \
           -v $(pwd):/go/src/jhqc.com/songcf/scene/ \
           139.198.2.55/soalib/golang:1.8 \
           /bin/bash "cd /go/src/jhqc.com/songcf/scene; make"


docker rmi scene:v1 --force
docker build -t="scene:v1" $(pwd)/scene

