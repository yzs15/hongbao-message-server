set -x

PERIOD=100

sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/all-test/bj-cycle.json
sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/all-test/nj-cycle.json

#ssh hbnj1 "
#sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/net.json ;
#sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/spb.json ;
#bash /root/thing/src-sync.sh ;
#" &
ssh hbnj1 "bash /root/thing/ms-stop.sh" &

bash scripts/restart-k8s-svs.sh lab9 sudo &
bash scripts/restart-k8s-svs.sh hbnj4 &

bash scripts/update-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json no_pull &
bash scripts/update-msd-cfg.sh hbnj4 configs/msd/bjnj/net-nj.json no_pull &

bash scripts/update-msd-cfg.sh lab3n configs/msd/bjnj/spb-bj.json no_pull &
bash scripts/update-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json no_pull &
wait

sleep 10

bash ../hongbao-log/scripts/update-log.sh hbnj4 &
ssh hbnj1 "bash /root/thing/ms-start.sh" &
bash scripts/update-thing.sh lab9 configs/things/all-test/bj-cycle.json 9 40 no_pull &
bash scripts/update-thing.sh hbnj2 configs/things/all-test/nj-cycle.json 2 8 no_pull &
#bash scripts/update-thing.sh hbnj4 configs/things/all-test/nj-cycle.json 4 8 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/all-test/nj-cycle.json 5 6 no_pull &
wait