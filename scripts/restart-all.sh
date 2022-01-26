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

LOG_ANA_SERVER=hbnj4

if [ "net" = "$ENV" ]; then # 互联网
  bash scripts/restart-k8s-svs.sh $NET_BJ_SERVER sudo &
  bash scripts/restart-k8s-svs.sh $NET_NJ_SERVER &
  wait

  # bash scripts/deploy-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json
  bash scripts/update-msd-cfg.sh $NET_BJ_SERVER configs/msd/bjnj/net-bj.json &
  bash scripts/update-msd-cfg.sh $NET_NJ_SERVER configs/msd/bjnj/net-nj.json &

elif [ "spb" = "$ENV" ]; then # 信息高铁
  # bash scripts/deploy-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json
  bash scripts/update-msd-cfg.sh $SPB_BJ_SERVER configs/msd/bjnj/spb-bj.json &
  bash scripts/update-msd-cfg.sh $SPB_NJ_SERVER configs/msd/bjnj/spb-nj.json &

else
  bash scripts/restart-k8s-svs.sh $NET_BJ_SERVER sudo &
  bash scripts/restart-k8s-svs.sh $NET_NJ_SERVER &
  wait

  # bash scripts/deploy-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json
  bash scripts/update-msd-cfg.sh $NET_BJ_SERVER configs/msd/bjnj/net-bj.json &
  bash scripts/update-msd-cfg.sh $NET_NJ_SERVER configs/msd/bjnj/net-nj.json &

  # bash scripts/deploy-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json
  bash scripts/update-msd-cfg.sh $SPB_BJ_SERVER configs/msd/bjnj/spb-bj.json &
  bash scripts/update-msd-cfg.sh $SPB_NJ_SERVER configs/msd/bjnj/spb-nj.json &
fi
wait

sleep 30

bash ../hongbao-log/scripts/update-log.sh $LOG_ANA_SERVER &
bash scripts/update-thing.sh lab9 configs/things/bjnj/bj-cycle.json 9 100 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/bjnj/nj-cycle.json 5 33 no_pull &
bash scripts/update-thing.sh hbnj4 configs/things/bjnj/nj-cycle.json 4 33 no_pull &
bash scripts/update-thing.sh hbnj2 configs/things/bjnj/nj-cycle.json 2 34 no_pull &
wait