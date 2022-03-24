ENV=net

cp configs/k8s/numrecd-each-node-detail.yaml.template configs/k8s/numrecd-detail.yaml.template

sed -i '' "s/[\#]\{0,\}. scripts\/conf-2C4C.sh/\#. scripts\/conf-2C4C.sh/g" scripts/test-all.sh
sed -i '' "s/[\#]\{0,\}. scripts\/conf-noise.sh/. scripts\/conf-noise.sh/g" scripts/test-all.sh

sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=10m/g" scripts/conf-noise.sh
sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=20/g" scripts/conf-noise.sh
sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-noise.sh
sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=2/g" scripts/conf-noise.sh
sed -i '' "s/MAX_REPLICA_NUM=\(.*\)/MAX_REPLICA_NUM=2/g" scripts/conf-noise.sh


# 没有噪音的不设限
LOG_PAR_DIR="/Volumes/Elements/logs-k8s-exp1-no-limit-no-noise-batch"
sed -i '' "s/WITH_NOISE=\(.*\)/WITH_NOISE=FALSE/g" scripts/conf-noise.sh
sed -i '' "s/NOISE_NUM_PER=\(.*\)/NOISE_NUM_PER=0/g" scripts/conf-noise.sh

for ((i=0;i<3;i++))
do
    for period in 15 25 50 100
    do
        echo "======+++++    start large $period $i    +++++========="
        while :
        do
            timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV large $period "$LOG_PAR_DIR"
            if [ $? -eq 0 ]; then
                break
            fi
        done
        echo "======+++++    end   large $period $i    +++++========="
    done
done

# 带有噪音的不设限
LOG_PAR_DIR="/Volumes/Elements/logs-k8s-exp1-no-limit-with-noise-batch"
sed -i '' "s/WITH_NOISE=\(.*\)/WITH_NOISE=TRUE/g" scripts/conf-noise.sh
sed -i '' "s/NOISE_NUM_PER=\(.*\)/NOISE_NUM_PER=1/g" scripts/conf-noise.sh

for ((i=0;i<3;i++))
do
    for period in 15 25 50 100
    do
        echo "======+++++    start large $period $i    +++++========="
        while :
        do
            timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV large $period "$LOG_PAR_DIR"
            if [ $? -eq 0 ]; then
                break
            fi
        done
        echo "======+++++    end   large $period $i    +++++========="
    done
done