#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: run.sh NAME ENV"
  exit 1
fi

NAME=$1
ENV=$2

HOST_IP=$(ifconfig docker0 | grep inet | awk 'NR==1{print $2}')

if [ "net" = $ENV ]; then
  WNAG_ZMQ_PORT=5557
else
  WNAG_ZMQ_PORT=5558
fi

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
docker rmi $(docker image ls -f dangling=true -q)
if [ "wang" = $NAME ]; then
  docker stop wangmsd
  docker rm wangmsd
  docker run -it \
    --name wangmsd \
    -p 8080:8080 \
    -p $WNAG_ZMQ_PORT:$WNAG_ZMQ_PORT \
    -p 5552:5552 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd -wang \
      -ws   0.0.0.0:8080 \
      -zmq  tcp://0.0.0.0:$WNAG_ZMQ_PORT \
      -log  0.0.0.0:5552 \
      -$ENV \
      -tend tcp://$HOST_IP:5543


elif [ "thing" = $NAME ]; then
  docker stop thingmsd
  docker rm thingmsd
  docker run -it \
    --name thingmsd \
    -p 5544:5544 \
    -p 5543:5543 \
    -p 5542:5542 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd -thing \
      -ws   0.0.0.0:5544 \
      -zmq  tcp://0.0.0.0:5543 \
      -log  0.0.0.0:5542 \
      -$ENV \
      -wend tcp://$HOST_IP:$WNAG_ZMQ_PORT \
      -kend 172.16.32.13 \
      -kend 172.16.32.14 \
      -kend 172.16.32.15
fi