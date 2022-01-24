bash scripts/deploy-msd.sh hbnj5 spb nj
bash scripts/update-msd.sh hbnj3 spb bj
bash scripts/update-thing.sh hbnj3 configs/things/spb-bj-cycle.json 3 2 &
bash scripts/update-thing.sh hbnj5 configs/things/spb-nj-cycle.json 5 3 &
wait

bash scripts/restart-k8s-svs.sh kbnj1 &
bash scripts/restart-k8s-svs.sh kbnj4 &
bash scripts/update-msd.sh kbnj3 net bj &
bash scripts/update-msd.sh kbnj5 net nj &

bash scripts/update-thing.sh kbnj2 configs/things/net-bj-cycle.json 2 3 &
bash scripts/update-thing.sh kbnj5 configs/things/net-nj-cycle.json 5 2 &




bash scripts/rm-k8s-log.sh lab3n &
bash scripts/rm-k8s-log.sh lab9 &
bash scripts/rm-k8s-log.sh hbnj4 &
bash scripts/rm-k8s-log.sh hbnj5 &
bash scripts/restart-k8s-svs.sh lab9 sudo &
bash scripts/restart-k8s-svs.sh hbnj4 &

bash scripts/deploy-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json
bash scripts/update-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json &
bash scripts/update-msd-cfg.sh hbnj4 configs/msd/bjnj/net-nj.json &

bash scripts/deploy-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json
bash scripts/update-msd-cfg.sh lab3n configs/msd/bjnj/spb-bj.json &
bash scripts/update-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json &

# LFZ 重启物

bash ../hongbao-log/scripts/deploy-log.sh lab9 &
bash ../hongbao-log/scripts/update-log.sh lab9 &

bash scripts/update-thing.sh lab9 configs/things/bjnj/bj-cycle.json 9 5 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/bjnj/nj-cycle.json 5 15 no_pull &

bash scripts/update-thing.sh lab9 configs/things/bjnj/bj-cycle.json 9 1 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/bjnj/nj-cycle.json 5 1 no_pull &
