#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: docker-run-msd.sh CONFIG [NO_PULL]"
  exit 1
fi

CONFIG=$1
NO_PULL=$2

PRO_DIR=$HOME/projects/hongbao-ms

ZMQ_PORT=8081
WS_PORT=8082

if [ -z $NO_PULL ]; then
 docker pull registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd
fi

# docker rmi $(docker image ls -f dangling=true -q)
docker stop msd
docker rm msd

if [[ $CONFIG =~ "net" ]]; then
    docker run -it \
      --name msd \
      -v $HOME/projects/hongbao-ms:/hongbao-ms \
      --device=/dev/ptp0 \
      --device=/dev/ptp1 \
      --network host \
      --cpu-period=100000 --cpu-quota=10000000 \
      registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
      bash -c "ulimit -n 65535 && msd -msdcfg $CONFIG"

else
    docker run -it \
      --name msd \
      -v $HOME/projects/hongbao-ms:/hongbao-ms \
      --device=/dev/ptp0 \
      --device=/dev/ptp1 \
      --network host \
      --cpuset-cpus 0-14 \
      registry.cn-beijing.aliyuncs.com/zhengsj/hongbao:msd \
      bash -c "ulimit -n 65535 && msd -msdcfg $CONFIG"
fi

# --cpu-period=100000 --cpu-quota=10000000 \
# --cpuset-cpus 0-14 \