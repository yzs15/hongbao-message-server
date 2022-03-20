ENV=net

cp configs/k8s/numrecd-normal-detail.yaml.template configs/k8s/numrecd-detail.yaml.template

LOG_PAR_DIR="/Volumes/Elements/logs-k8s-exp2-resource-limit-with-noise-batch"

CPU_REQUESTS=(100m 1 2 4 10m)
CPU_LIMITS=(200m 2 4 8 20)

sed -i '' "s/[\#]\{0,\}. scripts\/conf-2C4C.sh/\#. scripts\/conf-2C4C.sh/g" scripts/test-all.sh
sed -i '' "s/[\#]\{0,\}. scripts\/conf-noise.sh/. scripts\/conf-noise.sh/g" scripts/test-all.sh

for ((ci=0;ci<5;ci++))
do
    sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=${CPU_REQUESTS[ci]}/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=${CPU_LIMITS[ci]}/g" scripts/conf-noise.sh
    sed -i '' "s/CPU_SCALE_THRESHOLD=\(.*\)/CPU_SCALE_THRESHOLD=60/g" scripts/conf-noise.sh
    sed -i '' "s/INIT_REPLICA_NUM=\(.*\)/INIT_REPLICA_NUM=1/g" scripts/conf-noise.sh
    sed -i '' "s/MAX_REPLICA_NUM=\(.*\)/MAX_REPLICA_NUM=1000/g" scripts/conf-noise.sh

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
done