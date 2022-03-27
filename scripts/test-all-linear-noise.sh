. scripts/utils.sh

unset http_proxy
unset https_proxy
unset ALL_PROXY

if [ $# -lt 1 ]; then
    echo "usage: test-all.sh ENV SIZE PERIOD [LOG_PAR_DIR]"
    exit 1
fi

ENV=${1-all}
SIZE=${2-large}
export PERIOD=${3-"25"}
MODE=${5-"cycle"}
CUR_DATETIME=$(date +%m%d%H%M%S)

if [ "net" = "$ENV" ]; then
    THING_CONF_DIR="net-test"
elif [ "spb" = "$ENV" ]; then
    THING_CONF_DIR="spb-test"
elif [ "all" = "$ENV" ]; then
    THING_CONF_DIR="all-test"
else
    echo "no such env: " "$ENV"
    exit 1
fi



echo    "==================================="
echo    "===        configuration        ==="
echo    "==================================="


. scripts/conf-noise.sh

if [ "large" = "$SIZE" ]; then
    THING_NUM=64
elif [ "small" = "$SIZE" ]; then
    THING_NUM=4
else
    echo "no such size: " $SIZE
    exit 1
fi

if [ "linear" = "$MODE" ]; then
    TOTAL_TIME=$(((MAX_TIME_PER_SEC / LINEAR_RATIO) *2 + 10))
    UP_TIME=$((MAX_TIME_PER_SEC / LINEAR_RATIO))
    export TOTAL=$((THING_NUM * ((UP_TIME * LINEAR_RATIO + LINEAR_RATIO) * UP_TIME + 10 * MAX_TIME_PER_SEC )))
else
    export TOTAL=$((THING_NUM * (TOTAL_TIME * 1000 / PERIOD)))
fi



LOG_PAR_DIR=${4-"/Volumes/Elements/logs"}
LOG_DIR="$LOG_PAR_DIR/$CUR_DATETIME-$ENV-$TOTAL-$MAX_TIME_PER_SEC-$LINEAR_RATIO"
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    LOG_DIR=$LOG_DIR"-${CPU_REQUEST}C${CPU_LIMIT}C${CPU_SCALE_THRESHOLD}%-${INIT_REPLICA_NUM}"
fi
mkdir -p $LOG_DIR


CONF_CSV="$CUR_DATETIME,$ENV,$THING_NUM,$PERIOD,$TOTAL,$CPU_REQUEST,$CPU_LIMIT,$CPU_SCALE_THRESHOLD,$CPU_SCALE_UP_LIMIT,$INIT_REPLICA_NUM,$MAX_REPLICA_NUM"
echo $CONF_CSV > $LOG_DIR/config.csv
echo $CONF_CSV >> $LOG_PAR_DIR/config.csv

echo "THING_NUM:  " $THING_NUM
echo "PERIOD:     " $PERIOD
echo "TOTAL:      " $TOTAL
echo "TOTAL_TIME: " $TOTAL_TIME
echo "NOISE:      " $WITH_NOISE
echo "CPU_SCALE_THRESHOLD" $CPU_SCALE_THRESHOLD
echo "CPU_REQUEST " $CPU_REQUEST
echo "CPU_LIMIT" $CPU_LIMIT


##### sync scripts and other things to server
#for svr in lab3n lab9 hbnj1 hbnj2 hbnj4 hbnj5
#do
#    ensure_ok rsync -a ./* $svr:projects/hongbao-ms/  --exclude-from=.gitignore --exclude=data &
#done
#wait

echo    "==================================="
echo    "===  stop all components        ==="
echo    "==================================="

##### delete k8s services #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    {
        bash scripts/stop-k8s-svs.sh lab9  numrecd sudo
        bash scripts/stop-k8s-svs.sh lab9  numrecd-noise sudo
    } &
    {
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd-noise
    } &
fi

##### stop task switching #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    {
        TS_STOP=$(curl http://10.2.5.199:5555/stop)
        if [ "OK" != "$TS_STOP" ]; then
            exit 1
        fi
    } &
fi

##### stop all message server #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab3n "tmux send-keys -t msd:0.0 C-c" &
    ensure_ok ssh hbnj1 '"tmux send-keys -t msd:0.0 C-c"' &
fi
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab9  "tmux send-keys -t msd:0.0 C-c" &
    ensure_ok ssh hbnj4 '"tmux send-keys -t msd:0.0 C-c"' &
fi

##### stop all things #####
ssh lab3n "tmux kill-session -t thing" &
ssh lab9  "tmux kill-session -t thing" &
ensure_ok ssh hbnj1 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj2 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj4 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj5 '"tmux kill-session -t thing || : "' &
# ssh hbnj1 "bash /root/thing/ms-stop.sh" &
wait


echo    "==================================="
echo    "===  rename all logs            ==="
echo    "==================================="

##### rename k8s pod cpu usage log #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab3n "mv projects/occupy/records projects/occupy/records-$CUR_DATETIME" &
    ssh lab9  "mv projects/occupy/records projects/occupy/records-$CUR_DATETIME" &
    ssh hbnj4 "mv projects/occupy/records projects/occupy/records-$CUR_DATETIME" &
    ssh hbnj5 "mv projects/occupy/records projects/occupy/records-$CUR_DATETIME" &
fi

##### rename spb cpu usage log #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    log_path=$LOG_DIR/ts-cpu/
    ssh lab3n "cd /home/yuzishu/taskswitching/BJ_M1 && sudo mv worker_log worker_log.$CUR_DATETIME || : " &
    ssh lab9  "cd /home/yuzishu/taskswitching/BJ_M2 && sudo mv worker_log worker_log.$CUR_DATETIME || : " &
    ensure_ok ssh hbnj1 "'""cd /home/yuzishu/taskswitching/NJ_M1 && mv worker_log worker_log.$CUR_DATETIME || : ""'"  &
    ensure_ok ssh hbnj2 "'""cd /home/yuzishu/taskswitching/NJ_M2 && mv worker_log worker_log.$CUR_DATETIME || : ""'"  &
fi

##### remove k8s logs #####
ssh lab9  "sudo rm -rf /var/log/hongbao/*" &
ssh lab3n "sudo rm -rf /var/log/hongbao/*" &

##### rename message server log #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab3n "mv projects/hongbao-ms/msd.log projects/hongbao-ms/msd.log.$CUR_DATETIME" &
    ssh hbnj1 "mv projects/hongbao-ms/msd.log projects/hongbao-ms/msd.log.$CUR_DATETIME" &
fi
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab9  "mv projects/hongbao-ms/msd.log projects/hongbao-ms/msd.log.$CUR_DATETIME" &
    ssh hbnj4 "mv projects/hongbao-ms/msd.log projects/hongbao-ms/msd.log.$CUR_DATETIME" &
fi
wait



echo    "==================================="
echo    "===        reconfigure          ==="
echo    "==================================="

##### reconfigure things #####
sed -i '' "s/\"Period\":\(.*\)/\"Period\":\"${PERIOD}ms\",/g"   configs/things/$THING_CONF_DIR/bj-${MODE}.json
sed -i '' "s/\"Period\":\(.*\)/\"Period\":\"${PERIOD}ms\",/g"   configs/things/$THING_CONF_DIR/nj-${MODE}.json
sed -i '' "s/\"TotalTime\":\(.*\)/\"TotalTime\":\"${TOTAL_TIME}s\",/g" configs/things/$THING_CONF_DIR/bj-${MODE}.json
sed -i '' "s/\"TotalTime\":\(.*\)/\"TotalTime\":\"${TOTAL_TIME}s\",/g" configs/things/$THING_CONF_DIR/nj-${MODE}.json
sed -i '' "s/\"LoadNumPer\":\(.*\)/\"LoadNumPer\": ${LOAD_NUM_PER},/g"  configs/things/$THING_CONF_DIR/bj-${MODE}.json
sed -i '' "s/\"LoadNumPer\":\(.*\)/\"LoadNumPer\": ${LOAD_NUM_PER},/g"  configs/things/$THING_CONF_DIR/nj-${MODE}.json
sed -i '' "s/\"NoisNumPer\":\(.*\)/\"NoisNumPer\": ${NOISE_NUM_PER},/g" configs/things/$THING_CONF_DIR/bj-${MODE}.json
sed -i '' "s/\"NoisNumPer\":\(.*\)/\"NoisNumPer\": ${NOISE_NUM_PER},/g" configs/things/$THING_CONF_DIR/nj-${MODE}.json
if [[ "linear" == $MODE ]]; then
    sed -i '' "s/\"LinearRatio\":\(.*\)/\"LinearRatio\": ${LINEAR_RATIO},/g" configs/things/$THING_CONF_DIR/bj-${MODE}.json
    sed -i '' "s/\"LinearRatio\":\(.*\)/\"LinearRatio\": ${LINEAR_RATIO},/g" configs/things/$THING_CONF_DIR/nj-${MODE}.json
    sed -i '' "s/\"MaxTaskPerSec\":\(.*\)/\"MaxTaskPerSec\": ${MAX_TIME_PER_SEC},/g" configs/things/$THING_CONF_DIR/bj-${MODE}.json
    sed -i '' "s/\"MaxTaskPerSec\":\(.*\)/\"MaxTaskPerSec\": ${MAX_TIME_PER_SEC},/g" configs/things/$THING_CONF_DIR/nj-${MODE}.json
fi
##### reconfigure real things #####
# ssh hbnj1 "
# sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/net.json ;
# sed -i \"s/[0-9]\{1,\}ms/${PERIOD}ms/g\" /root/numrec-thing/configs/test/spb.json ;
# bash /root/thing/src-sync.sh ;
# " &

##### reconfigure log analyser #####
sed -i '' "s/TOTAL = [0-9]\{1,\}/TOTAL = ${TOTAL}/g" ../hongbao-log/src/analyzer.py

##### reconfigure K8s config #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    bash scripts/build-yml.sh
    scp configs/k8s/numrecd.yml lab9:projects/k8s/numrecd.yml  &
    scp configs/k8s/numrecd.yml hbnj4:projects/k8s/numrecd.yml &
fi
wait



echo    "==================================="
echo    "===        restart platform     ==="
echo    "==================================="

##### restart Kubernetes #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    if [ "TRUE" = $WITH_NOISE ]; then   
        {
            scp configs/k8s/numrecd-noise.yml lab9:projects/k8s/numrecd-noise.yml
            bash scripts/restart-k8s-svs.sh lab9 numrecd-noise sudo
        } &
        {
            scp configs/k8s/numrecd-noise.yml hbnj4:projects/k8s/numrecd-noise.yml
            bash scripts/restart-k8s-svs.sh hbnj4 numrecd-noise
        } &
    fi
    wait
    bash scripts/restart-k8s-svs.sh lab9  numrecd sudo &
    bash scripts/restart-k8s-svs.sh hbnj4 numrecd &
    
fi

##### restart Task switching #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    TS_RESTART=$(curl http://10.2.5.199:5555/start)
    if [ "OK" != "$TS_RESTART" ]; then
        exit 1
    fi
fi
wait



echo    "==================================="
echo    "===     start message server    ==="
echo    "==================================="

##### start all message server #####
NO_PULL="no"
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    bash scripts/update-msd-cfg.sh lab3n configs/msd/bjnj/spb-bj.json $NO_PULL &
    bash scripts/update-msd-cfg.sh hbnj1 configs/msd/bjnj/spb-nj.json $NO_PULL &
fi
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    bash scripts/update-msd-cfg.sh lab9  configs/msd/bjnj/net-bj.json $NO_PULL &
    bash scripts/update-msd-cfg.sh hbnj4 configs/msd/bjnj/net-nj.json $NO_PULL &
fi
wait

echo "sleep 20s"
sleep 20



echo    "==================================="
echo    "=== start things & log analyzer ==="
echo    "==================================="

##### start log analyzer #####
bash ../hongbao-log/scripts/update-log.sh hbnj4 &

##### start all things #####
# ssh hbnj1 "bash /root/thing/ms-start.sh" &

if [ "large" = "$SIZE" ]; then
    bash scripts/update-thing.sh lab9  configs/things/$THING_CONF_DIR/bj-${MODE}.json 9 40 no_pull &
    bash scripts/update-thing.sh hbnj2 configs/things/$THING_CONF_DIR/nj-${MODE}.json 2 8  no_pull &
    bash scripts/update-thing.sh hbnj5 configs/things/$THING_CONF_DIR/nj-${MODE}.json 5 16 no_pull &
elif [ "small" = "$SIZE" ]; then
    bash scripts/update-thing.sh lab9 configs/things/$THING_CONF_DIR/bj-${MODE}.json  9 2 no_pull &
    bash scripts/update-thing.sh hbnj5 configs/things/$THING_CONF_DIR/nj-${MODE}.json 5 2 no_pull &
fi

# bash scripts/update-thing.sh hbnj1 configs/things/$THING_CONF_DIR/nj-${MODE}.json 1 40 no_pull &
# bash scripts/update-thing.sh hbnj2 configs/things/$THING_CONF_DIR/nj-${MODE}.json 2 40 no_pull &
# bash scripts/update-thing.sh hbnj4 configs/things/$THING_CONF_DIR/nj-${MODE}.json 4 40 no_pull &
# bash scripts/update-thing.sh hbnj5 configs/things/$THING_CONF_DIR/nj-${MODE}.json 5 40 no_pull &
wait

echo "sleep 20s"
sleep 20



# echo    "==================================="
# echo    "===  start monitor              ==="
# echo    "==================================="

# if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
#     bash ../machine-monitor/scripts/deploy.sh lab3n &
#     bash ../machine-monitor/scripts/deploy.sh lab9  &
#     bash ../machine-monitor/scripts/deploy.sh hbnj4 &
#     bash ../machine-monitor/scripts/deploy.sh hbnj5 &
# fi
# wait



echo    "==================================="
echo    "===  start testing              ==="
echo    "==================================="

##### start testing #####
cd ../hongbao-log
source ./venv/bin/activate

if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    python3 ./src/mock_wang.py spb $((TOTAL_TIME+60))
    echo "sleep $((80+TOTAL_TIME))s"
    sleep $((80+TOTAL_TIME))
fi
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    python3 ./src/mock_wang.py net $((TOTAL_TIME+60))
    echo "sleep $((80+TOTAL_TIME))s"
    sleep $((80+TOTAL_TIME))
fi


echo    "==================================="
echo    "===  stop monitor              ==="
echo    "==================================="

if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab3n "tmux send-keys -t ctnm:0.0 C-c C-m" &
    ssh lab9  "tmux send-keys -t ctnm:0.0 C-c C-m" &
    ensure_ok ssh hbnj4 '"tmux send-keys -t ctnm:0.0 C-c C-m"' &
    ensure_ok ssh hbnj5 '"tmux send-keys -t ctnm:0.0 C-c C-m"' &
fi
wait


echo    "==================================="
echo    "===  stop all components        ==="
echo    "==================================="

##### delete k8s services #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    {
        bash scripts/stop-k8s-svs.sh lab9  numrecd sudo
        bash scripts/stop-k8s-svs.sh lab9  numrecd-noise sudo
    } &
    {
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd-noise
    } &
fi

##### stop task switching #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    {
        TS_STOP=$(curl http://10.2.5.199:5555/stop)
        if [ "OK" != "$TS_STOP" ]; then
            exit 1
        fi
    } &
fi

##### stop all message server #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab3n "tmux send-keys -t msd:0.0 C-c" &
    ensure_ok ssh hbnj1 '"tmux send-keys -t msd:0.0 C-c"' &
fi
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    ssh lab9  "tmux send-keys -t msd:0.0 C-c" &
    ensure_ok ssh hbnj4 '"tmux send-keys -t msd:0.0 C-c"' &
fi

##### stop all things #####
ssh lab3n "tmux kill-session -t thing" &
ssh lab9  "tmux kill-session -t thing" &
ensure_ok ssh hbnj1 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj2 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj4 '"tmux kill-session -t thing || : "' &
ensure_ok ssh hbnj5 '"tmux kill-session -t thing || : "' &
# ssh hbnj1 "bash /root/thing/ms-stop.sh" &
wait


echo    "==================================="
echo    "===  download logs              ==="
echo    "==================================="

cd ../hongbao-ms

##### download k8s pod cpu usage log #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    for svr in lab3n lab9 hbnj4 hbnj5; do
        log_path=$LOG_DIR/k8s-cpu/$svr
        mkdir -p "$log_path"
        ensure_ok rsync -amP $svr:"projects/occupy/records/*" "$log_path" &
    done
fi

##### download spb cpu usage log #####
if [[ "spb" == $ENV ]] || [[ "all" == $ENV ]]; then
    log_path=$LOG_DIR/ts-cpu/
    mkdir -p $log_path/BJ_M1
    mkdir -p $log_path/BJ_M2
    mkdir -p $log_path/NJ_M1
    mkdir -p $log_path/NJ_M2
    scp -r lab3n:/home/yuzishu/taskswitching/BJ_M1/worker_log/*.log $log_path/BJ_M1 &
    scp -r lab9:/home/yuzishu/taskswitching/BJ_M2/worker_log/*.log  $log_path/BJ_M2 &
    ensure_ok scp -r hbnj1:/home/yuzishu/taskswitching/NJ_M1/worker_log/*.log $log_path/NJ_M1 &
    ensure_ok scp -r hbnj2:/home/yuzishu/taskswitching/NJ_M2/worker_log/*.log $log_path/NJ_M2 &
fi

##### download log analyser log #####
{
    remote_log_dir=$(ssh hbnj4 "ls -lt ~/projects/hongbao-log/logs | grep \"^d\" | grep 2022"  | head -n 1 | awk '{print $9}')
    rsync -azP hbnj4:projects/hongbao-log/logs/$remote_log_dir $LOG_DIR/
    ana_result=$(cat $LOG_DIR/$remote_log_dir/result.csv)
    echo "$CUR_DATETIME,""$ana_result" >> $LOG_PAR_DIR/result.csv
} &
wait

echo    "==================================="
echo    "===  analyze CPU usage          ==="
echo    "==================================="
bash scripts/analyze.sh $LOG_DIR &


##### delete k8s services #####
if [[ "net" == $ENV ]] || [[ "all" == $ENV ]]; then
    echo    "==================================="
    echo    "===  delete K8s services        ==="
    echo    "==================================="

    {
        bash scripts/stop-k8s-svs.sh lab9  numrecd sudo
        bash scripts/stop-k8s-svs.sh lab9  numrecd-noise sudo
    } &
    {
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd
        bash scripts/stop-k8s-svs.sh hbnj4 numrecd-noise
    } &
fi
wait