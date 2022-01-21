#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: docker-run-msd.sh ENV LOC"
  exit 1
fi

ENV=$1
LOC=$2

PRO_DIR=/root/projects/hongbao-ms

if [ "net" = $ENV ]; then
  ZMQ_PORT=8081
  WS_PORT=8082
  if [ "bj" = $LOC ]; then
    KENDS="-kend 172.16.32.12"
    NSEND="172.16.32.13:8084"
    ZMQ_OUT="tcp://172.16.32.14:8081"
#    KENDS="-kend 10.208.104.9"
#    NSEND="58.213.121.2:10027"
#    ZMQ_OUT="tcp://159.226.41.229:7102"
    MAC="02:42:ac:11:00:01"
  else
    KENDS="-kend 172.16.32.15"
    NSEND="172.16.32.13:8084"
    ZMQ_OUT="tcp://172.16.32.16:8081"
#    KENDS="-kend 172.16.32.13 -kend 172.16.32.14 -kend 172.16.32.15"
#    NSEND="172.16.32.13:8080"
#    ZMQ_OUT="tcp://58.213.121.2:10025"
    MAC="02:42:ac:11:00:02"
  fi
else
  ZMQ_PORT=8081
  WS_PORT=8082
  KENDS=""
  if [ "bj" = $LOC ]; then
    NSEND="192.168.143.2:8084"
    ZMQ_OUT="tcp://192.168.143.3:8081"
    MAC="02:42:ac:11:00:01"
  else
    NSEND="192.168.143.2:8084"
    ZMQ_OUT="tcp://192.168.143.5:8081"
    MAC="02:42:ac:11:00:02"
  fi
fi

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
# docker rmi $(docker image ls -f dangling=true -q)
docker stop msd
docker rm msd
docker run -it \
  --name msd \
  -p $WS_PORT:5554 \
  -p $ZMQ_PORT:5553 \
  -p 5552:5552 \
  -v /root/projects/hongbao-ms:/hongbao-ms \
  --mac-address $MAC \
  registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
  msd \
    -ws   0.0.0.0:5554 \
    -zmq  tcp://0.0.0.0:5553 \
    -log-path /hongbao-ms \
    -nsend $NSEND \
    -zmq-out $ZMQ_OUT \
    -$ENV $KENDS