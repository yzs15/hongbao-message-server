#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: restart-all.sh ENV"
  exit 1
fi

ENV=$1

SPB_BJ_SERVER=lab3n
SPB_NJ_SERVER=hbnj1

NET_BJ_SERVER=lab9
NET_NJ_SERVER=hbnj4
KUBE_SERVER=(lab3n lab9 hbnj4 hbnj5)

LOG_ANA_SERVER=lab9

if [ "net" = "$ENV" ]; then # 互联网
  for svr in "${KUBE_SERVER[@]}"
  do
    echo $svr
    bash scripts/rm-k8s-log.sh "$svr" &
  done
  wait
  sleep 3

  bash scripts/restart-k8s-svs.sh $NET_BJ_SERVER sudo &
  bash scripts/restart-k8s-svs.sh $NET_NJ_SERVER &

  # bash scripts/deploy-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json
  bash scripts/update-msd-cfg.sh $NET_BJ_SERVER configs/msd/bjnj/net-bj.json &
  bash scripts/update-msd-cfg.sh $NET_NJ_SERVER configs/msd/bjnj/net-nj.json &

else # 信息高铁
  # bash scripts/deploy-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json
  bash scripts/update-msd-cfg.sh $SPB_BJ_SERVER configs/msd/bjnj/spb-bj.json &
  bash scripts/update-msd-cfg.sh $SPB_NJ_SERVER configs/msd/bjnj/spb-nj.json &
fi

wait
sleep 30

# bash ../hongbao-log/scripts/update-log.sh $LOG_ANA_SERVER &

# bash scripts/update-thing.sh hbnj5 configs/things/bjnj/nj-cycle.json 5 25 &
bash scripts/update-thing.sh hbnj5 configs/things/spb-test/nj-cycle.json 5 25 &
wait