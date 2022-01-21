#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: docker-run-msd.sh CONFIG"
  exit 1
fi

CONFIG=$1

PRO_DIR=$HOME/projects/hongbao-ms

ZMQ_PORT=8081
WS_PORT=8082

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
# docker rmi $(docker image ls -f dangling=true -q)
docker stop msd
docker rm msd
docker run -it \
  --name msd \
  -p $WS_PORT:5554 \
  -p $ZMQ_PORT:5553 \
  -p 5552:5552 \
  -v $HOME/projects/hongbao-ms:/hongbao-ms \
  registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
  msd -msdcfg $CONFIG