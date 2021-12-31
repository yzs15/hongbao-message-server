#!/bin/bash
cd $(dirname "$0")
cd ..

NAME=$1

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
if [ "wang" = $NAME ]; then
  docker stop wangmsd
  docker rm wangmsd
  docker run -it \
    --name wangmsd \
    -p 5552-5554:5552-5554 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd --wang \
      --ws   0.0.0.0:5554 \
      --zmq  tcp://0.0.0.0:5553 \
      --log  0.0.0.0:5552 \
      --tend tcp://127.0.0.1:5543


elif [ "thing" = $NAME ]; then
  docker stop thingmsd
  docker rm thingmsd
  docker run -it \
    --name thingmsd -p 5542-5544:5542-5544 \
    registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
    msd --thing \
      --ws   0.0.0.0:5544 \
      --zmq  tcp://0.0.0.0:5543 \
      --log  0.0.0.0:5542 \
      --wend tcp://127.0.0.1:5553 \
      --kend 172.16.32.12:32101 \
      --kend 172.16.32.13:32101 \
      --kend 172.16.32.14:32101 \
      --kend 172.16.32.15:32101

else
  echo "usage: run.sh NAME"
fi