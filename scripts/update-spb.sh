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


