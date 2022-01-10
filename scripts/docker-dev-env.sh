#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: run.sh ENV LOC"
  exit 1
fi

ENV=$1
LOC=$2

if [ "net" = $ENV ]; then
  ZMQ_PORT=5557
  WS_PORT=8080
  if [ "bj" = $LOC ]; then
    PRO_DIR="/home/zsj/projects/hongbao-ms"
    MAC="02:42:ac:11:00:01"
  else
    PRO_DIR="/root/projects/hongbao-ms"
    MAC="02:42:ac:11:00:02"
  fi
else
  ZMQ_PORT=5558
  WS_PORT=8080
  PRO_DIR="/root/projects/hongbao-ms"
fi

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
if [ "" != "$(docker image ls -f dangling=true -q)" ]; then
  docker rmi $(docker image ls -f dangling=true -q)
fi
docker stop msd-dev
docker rm msd-dev

docker run -it \
  --name msd-dev \
  -p $WS_PORT:5554 \
  -p $ZMQ_PORT:5553 \
  -p 5552:5552 \
  -v $PRO_DIR:/hongbao-ms \
  --mac-address $MAC \
  registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
  bash