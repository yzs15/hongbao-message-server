ENV=net
MODE=linear
# LOG_PAR_DIR="/Volumes/Elements/logs-yuzishu-3-24-valid-linear-with-noise"
LOG_PAR_DIR="/Volumes/Elements/logs-yuzishu-3-26-valid-linear-with-noise-18C-20C-1x2-limit"
# export MAX_TIME_PER_SEC=120
# export LINEAR_RATIO=4
# sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=10m/g" scripts/conf-noise.sh
# sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-noise.sh
# timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
# exit 


# test 4C8C
# sed -i '' "s/CPU_REQUEST=[0-9]\{1,\}/CPU_REQUEST=4/g" scripts/test-all.sh
# sed -i '' "s/CPU_LIMIT=[0-9]\{1,\}/CPU_LIMIT=8/g" scripts/test-all.sh
# TS_STOP=$(curl http://10.2.5.199:5555/stop)

# ENV=spb
# for ((i=60;i<=140;i=i+10))
# do 
#     export MAX_TIME_PER_SEC=$i
#     export LINEAR_RATIO=50
#     timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
# done

# TS_STOP=$(curl http://10.2.5.199:5555/stop)
# TS_STOP=$(curl http://10.2.5.199:5555/stop)
# TS_STOP=$(curl http://10.2.5.199:5555/stop)

# ENV=net
# for ((i=30;i<=100;i=i+10))
# do 
    # export MAX_TIME_PER_SEC=$i
    # export LINEAR_RATIO=1
    # cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
    # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=4/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=8/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=20/g" scripts/conf-noise.sh
    # timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

    # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=2/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=4/g" scripts/conf-noise.sh
    # timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

    # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=1/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=2/g" scripts/conf-noise.sh
    # timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE

    
    
    # cp configs/k8s/numrecd-each-node-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
    # sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=18/g" scripts/conf-noise.sh
    # sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-noise.sh
    # sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=2/g" scripts/conf-noise.sh
    # timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
    # sleep 60
# done

for ((i=60;i<=60;i=i+10))
do 
    export MAX_TIME_PER_SEC=$i
    export LINEAR_RATIO=10
    cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
    sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=10m/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-noise.sh
    timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
done

for ((i=20;i<=50;i=i+10))
do 
    export MAX_TIME_PER_SEC=$i
    export LINEAR_RATIO=1
    cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template
    sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=4/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=8/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=10/g" scripts/conf-noise.sh
    timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
done



# for ((i=0;i<3;i++))
# do
#     for period in 10 12 13 14 15 25 50 100 200 400 800 1600
#     do
#         echo "======+++++    start large $period $i    +++++========="
#         while :
#         do
#             timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV large $period "$LOG_PAR_DIR"
#             if [ $? -eq 0 ]; then
#                 break
#             fi
#         done
#         echo "======+++++    end   large $period $i    +++++========="
#     done

#     for period in 200 400 800 1600
#     do
#         echo "======+++++    start small $period $i    +++++========="
#         while :
#         do
#             timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV small $period "$LOG_PAR_DIR"
#             if [ $? -eq 0 ]; then
#                 break
#             fi
#         done
#         echo "======+++++    end   small $period $i    +++++========="
#     done
# done

# ENV=spb
# for ((i=50;i<=140;i=i+10))
# do 
#     export MAX_TIME_PER_SEC=$i
#     export LINEAR_RATIO=10
#     timeout --foreground "$((20*60))" bash scripts/test-all-linear-noise.sh $ENV large 10 "$LOG_PAR_DIR" $MODE
# done

# TS_STOP=$(curl http://10.2.5.199:5555/stop)
# TS_STOP=$(curl http://10.2.5.199:5555/stop)
# TS_STOP=$(curl http://10.2.5.199:5555/stop)
# bash scripts/yuzishu_test.sh