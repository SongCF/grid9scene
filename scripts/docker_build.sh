#!/bin/bash

make clean
# make clean-deps

docker login -usoa -pAabb0011 139.198.2.55
docker pull 139.198.2.55/soalib/golang:1.8

# 成功build后输出build success
# 失败后输出build failed，并exit 1
docker stop build_scene || echo "skip stop build_scene"
docker rm build_scene || echo "skip rm build_scene"
docker run -v $(pwd):/go/src/jhqc.com/songcf/scene/ \
           -w /go/src/jhqc.com/songcf/scene/ \
           --name build_scene \
           139.198.2.55/soalib/golang:1.8 \
           /bin/bash -c "./gopack get-deps && go build -o scene && echo \"\n---build success---\n\" || echo \"---build failed---\" && exit 1"

if [ $(docker wait build_scene) -eq 0 ];then exit 1; else echo ok; fi
docker stop build_scene
docker rm build_scene

docker rmi scene:build --force || echo "skip rmi scene:build"
docker build -t="scene:build" .
