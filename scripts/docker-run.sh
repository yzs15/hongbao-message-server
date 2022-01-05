#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: run.sh NAME"
  exit 1
fi

NAME=$1

HOST_IP=$(ifconfig docker0 | grep inet | awk 'NR==1{print $2}')

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
docker rmi $(docker image ls -f dangling=true -q)
if [ "wang" = $NAME ]; then
  docker stop wangmsd
  docker rm wangmsd
  docker run -it \
    --name wangmsd \
    -p 5552-5554:5552-5554 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd -wang \
      -ws   0.0.0.0:5554 \
      -zmq  tcp://0.0.0.0:5553 \
      -log  0.0.0.0:5552 \
      -net \
      -tend tcp://$HOST_IP:5543


elif [ "thing" = $NAME ]; then
  docker stop thingmsd
  docker rm thingmsd
  docker run -it \
    --name thingmsd -p 5542-5544:5542-5544 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd -thing \
      -ws   0.0.0.0:5544 \
      -zmq  tcp://0.0.0.0:5543 \
      -log  0.0.0.0:5542 \
      -net \
      -wend tcp://$HOST_IP:5553 \
      -kend 172.16.32.13 \
      -kend 172.16.32.14 \
      -kend 172.16.32.15
fi