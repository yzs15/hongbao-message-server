#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: docker-run-thing.sh CONFIG NODE IDX"
  exit 1
fi

CONFIG=$1
NODE=$(printf "%02x" $2)
IDX=$(printf "%02x" $3)

MAC_ADDR=02:42:ac:12:$NODE:$IDX

CONTAINER_NAME=thing-$IDX

docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
# docker rmi $(docker image ls -f dangling=true -q)
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker run -it \
  --name $CONTAINER_NAME \
  -v /root/projects/hongbao-ms:/hongbao-ms \
  --mac-address $MAC_ADDR \
  registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
  thingcli \
    -config /hongbao-ms/$CONFIG