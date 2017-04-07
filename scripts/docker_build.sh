#!/bin/bash

git checkout .

docker login -usoa -pAabb0011 139.198.2.55
docker pull 139.198.2.55/soalib/golang:1.8
docker run -v 

go install jhqc.com/songcf/scene

docker rmi scene_go:v1 --force
docker build -t="scene_go:v1" $GOPATH/src/jhqc.com/songcf/scene

