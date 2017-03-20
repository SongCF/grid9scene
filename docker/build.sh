#!/bin/bash

git checkout .
rm -rf log

# 清空本地依赖，重新获取
rm -rf deps/cluster
rm -rf deps/reloader
rm -rf deps/notify
rm -rf deps/super
rm -rf deps/framework
rm -rf deps/proto

chmod +x rebar
chmod +x ./docker/start.sh
make clean
make deps
cp ./docker/Dockerfile ./../Dockerfile
docker rmi 172.16.1.5:5000/plat_scene:v1 --force
docker build -t="172.16.1.5:5000/plat_scene:v1" ./../
rm -f ./../Dockerfile
