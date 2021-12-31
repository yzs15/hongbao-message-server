#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

IMAGE_NAME=registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd

docker stop wangmsd
docker rm wangmsd

docker stop thingmsd
docker rm thingmsd

docker rmi $IMAGE_NAME
docker build -t $IMAGE_NAME .
docker push $IMAGE_NAME

