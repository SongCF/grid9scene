#!/bin/bash

git checkout .

go install jhqc.com/songcf/scene

docker rmi scene_go:v1 --force
docker build -t="scene_go:v1" $GOPATH/src/jhqc.com/songcf/scene

