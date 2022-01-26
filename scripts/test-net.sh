set -x

PERIOD=100

sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/net-test/bj-cycle.json
sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/net-test/nj-cycle.json

ssh hbnj1 "
sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/net.json ;
sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/spb.json ;
bash /root/thing/src-sync.sh ;
" &

#bash scripts/rm-k8s-log.sh lab3n &
#bash scripts/rm-k8s-log.sh lab9 &
#bash scripts/rm-k8s-log.sh hbnj4 &
#bash scripts/rm-k8s-log.sh hbnj5 &
#wait
bash scripts/restart-k8s-svs.sh lab9 sudo &
bash scripts/restart-k8s-svs.sh hbnj4 &
wait

bash scripts/update-msd-cfg.sh lab9 configs/msd/bjnj/net-bj.json no_pull &
bash scripts/update-msd-cfg.sh hbnj4 configs/msd/bjnj/net-nj.json no_pull &
wait

sleep 20

bash ../hongbao-log/scripts/update-log.sh lab9 &
ssh hbnj1 "bash /root/thing/ms-start-test.sh" &
bash scripts/update-thing.sh lab9 configs/things/net-test/bj-cycle.json 9 5 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/net-test/nj-cycle.json 5 15 no_pull &
wait