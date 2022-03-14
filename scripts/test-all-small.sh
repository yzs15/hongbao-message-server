set -x

unset http_proxy
unset https_proxy
unset ALL_PROXY

PERIOD=800
THING_NUM=4
TOTAL=$((40*$THING_NUM*1000/PERIOD))

##### reconfigure things #####
sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/all-test/bj-cycle.json
sed -i '' "s/[0-9]\{1,\}ms/${PERIOD}ms/g" configs/things/all-test/nj-cycle.json

##### reconfigure real things #####
# ssh hbnj1 "
# sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/net.json ;
# sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/spb.json ;
# bash /root/thing/src-sync.sh ;
# " &


##### reconfigure log analyser #####
sed -i '' "s/TOTAL = [0-9]\{1,\}/TOTAL = ${TOTAL}/g" ../hongbao-log/src/analyzer.py

##### stop all message server #####
ssh lab3n "tmux send-keys -t msd:0.0 C-c" &
ssh hbnj1 "tmux send-keys -t msd:0.0 C-c" &
ssh lab9  "tmux send-keys -t msd:0.0 C-c" &
ssh hbnj4 "tmux send-keys -t msd:0.0 C-c" &

##### stop all things #####
ssh lab3n "tmux kill-session -t thing" &
ssh lab9  "tmux kill-session -t thing" &
ssh hbnj1 "tmux kill-session -t thing" &
ssh hbnj2 "tmux kill-session -t thing" &
ssh hbnj4 "tmux kill-session -t thing" &
ssh hbnj5 "tmux kill-session -t thing" &
#ssh hbnj1 "bash /root/thing/ms-stop.sh" &
wait

##### restart Kubernetes #####
ssh lab9  "sudo rm -rf /var/log/hongbao/*"
ssh lab3n "sudo rm -rf /var/log/hongbao/*"
bash scripts/restart-k8s-svs.sh lab9 sudo &
bash scripts/restart-k8s-svs.sh hbnj4 &

##### restart Task switching #####
TS_STOP=$(curl http://10.2.5.199:5555/stop)
if [ "OK" != "$TS_STOP" ]; then
  exit 1
fi
TS_RESTART=$(curl http://10.2.5.199:5555/start)
if [ "OK" != "$TS_RESTART" ]; then
  exit 1
fi
wait

##### start all message server #####
NO_PULL="no"
bash scripts/update-msd-cfg.sh lab3n configs/msd/bjnj/spb-bj.json $NO_PULL &
bash scripts/update-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json $NO_PULL &
bash scripts/update-msd-cfg.sh lab9  configs/msd/bjnj/net-bj.json $NO_PULL &
bash scripts/update-msd-cfg.sh hbnj4 configs/msd/bjnj/net-nj.json $NO_PULL &
wait

sleep 10

##### start log analyzer #####
bash ../hongbao-log/scripts/update-log.sh hbnj4 &

##### start all things #####
# ssh hbnj1 "bash /root/thing/ms-start.sh" &
#bash scripts/update-thing.sh lab9  configs/things/all-test/bj-cycle.json 9 40 no_pull &
#bash scripts/update-thing.sh hbnj2 configs/things/all-test/nj-cycle.json 2 8  no_pull &
#bash scripts/update-thing.sh hbnj5 configs/things/all-test/nj-cycle.json 5 16 no_pull &

bash scripts/update-thing.sh lab9  configs/things/all-test/bj-cycle.json 9 2 no_pull &
bash scripts/update-thing.sh hbnj5 configs/things/all-test/nj-cycle.json 5 2 no_pull &

# bash scripts/update-thing.sh hbnj1 configs/things/all-test/nj-cycle.json 1 40 no_pull &
# bash scripts/update-thing.sh hbnj2 configs/things/all-test/nj-cycle.json 2 40 no_pull &
# bash scripts/update-thing.sh hbnj4 configs/things/all-test/nj-cycle.json 4 40 no_pull &
# bash scripts/update-thing.sh hbnj5 configs/things/all-test/nj-cycle.json 5 40 no_pull &
wait

sleep 30

##### start testing #####
cd ../hongbao-log
source ./venv/bin/activate

python3 ./src/mock_wang.py spb
sleep $((30))

#python3 ./src/mock_wang.py net
#sleep $((30))

bash scripts/down-logs.sh hbnj4