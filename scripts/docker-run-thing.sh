#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: docker-run-thing-cfg.sh CONFIG NODE IDX [NO_PULL]"
  exit 1
fi

CONFIG=$1
NODE=$(printf "%02x" $2)
IDX=$(printf "%02x" $3)
NO_PULL=$4

MAC_ADDR=02:42:ac:12:$NODE:$IDX

CONTAINER_NAME=thing-$IDX

# if [ -z $NO_PULL ]; then
#   docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
# fi
# docker rmi $(docker image ls -f dangling=true -q)
# docker stop $CONTAINER_NAME
# docker rm $CONTAINER_NAME
docker run -it \
  --name $CONTAINER_NAME \
  -v $HOME/projects/hongbao-ms:/hongbao-ms \
  --device=/dev/ptp0 \
  --device=/dev/ptp1 \
  --mac-address $MAC_ADDR \
  registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
  thingcli \
    -config /hongbao-ms/$CONFIG