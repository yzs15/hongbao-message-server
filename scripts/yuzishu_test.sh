ENV=net
MODE=linear
LOG_PAR_DIR="/Volumes/Elements/logs-yuzishu-3-27-valid-linear-no-noise-k8s-limit-exp1"

# ENV=net
# for linear_ratio in 10 50 70 #1
# do
#     for ((max_task_ps=30;max_task_ps<180;max_task_ps=max_task_ps+20))
#     do
#         export MAX_TIME_PER_SEC=$max_task_ps
#         export LINEAR_RATIO=$linear_ratio
#         cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#         sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=4/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=8/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=10/g" scripts/conf-2C4C.sh
#         timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#         cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#         sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-2C4C.sh
#         timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#         cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#         sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=10m/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-2C4C.sh
#         sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-2C4C.sh
#         timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
#     done
# done

ENV=spb
for linear_ratio in 10 50 70 #1
do
    for ((max_task_ps=10;max_task_ps<180;max_task_ps=max_task_ps+20))
    do
        export MAX_TIME_PER_SEC=$max_task_ps
        export LINEAR_RATIO=$linear_ratio
        timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
    done
done
TS_STOP=$(curl http://10.2.5.199:5555/stop)
TS_STOP=$(curl http://10.2.5.199:5555/stop)
TS_STOP=$(curl http://10.2.5.199:5555/stop)

# ENV=net
# for ((i=1;i<=10;i++))
# do 
#     export MAX_TIME_PER_SEC=60
#     export LINEAR_RATIO=$i
#     cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=4/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=8/g" scripts/conf-2C4C.sh

#     # sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#     # sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=40/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#     sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=10/g" scripts/conf-2C4C.sh
#     timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE


#     cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=2/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=4/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=10/g" scripts/conf-2C4C.sh
#     timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#     # cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=10/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE


#     cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-2C4C.sh
#     sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=5/g" scripts/conf-2C4C.sh
#     timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#     # cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

#     # cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=10m/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
    
#     # cp configs/k8s/numrecd-each-node-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
#     # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=18/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-2C4C.sh
#     # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=2/g" scripts/conf-2C4C.sh
#     # timeout --foreground "$((20*60))" bash scripts/test-all-linear.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
# done

LOG_PAR_DIR="/Volumes/Elements/logs-yuzishu-3-27-valid-linear-no-noise-k8s-limit-exp2"
cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-2C4C.sh
sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-2C4C.sh
sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-2C4C.sh
sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-2C4C.sh

    
for period in 10 12 13 14 15 25 50 100 200 400 800 1600 # 8
do
    echo "======+++++    start large $period $i    +++++========="
    # while :
    # do
        ENV=net
        timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV large $period "$LOG_PAR_DIR"
        ENV=spb
        timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV large $period "$LOG_PAR_DIR"
    # done
    echo "======+++++    end   large $period $i    +++++========="
done

for period in 200 400 800 1600
do
    echo "======+++++    start small $period $i    +++++========="
    # while :
    # do
        ENV=net
        timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV small $period "$LOG_PAR_DIR"
        ENV=spb
        timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV small $period "$LOG_PAR_DIR"
    # done
    echo "======+++++    end   small $period $i    +++++========="
done



