#!/bin/bash

git checkout .
make clean
# TODO delete
make clean-deps

docker login -usoa -pAabb0011 139.198.2.55
docker pull 139.198.2.55/soalib/golang:1.8
docker run --rm \
           -v $(pwd):/go/src/jhqc.com/songcf/scene/ \
           -w /go/src/jhqc.com/songcf/scene/
           139.198.2.55/soalib/golang:1.8 \
           /bin/bash ./gopack get-deps & go build -o scene


docker rmi scene:v1 --force
docker build -t="scene:v1" .

