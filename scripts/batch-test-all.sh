ENV=net

LOG_PAR_DIR="/Volumes/Elements/logs-k8s-4C8C-batch-3-23-valid"
# sed -i '' "s/[\#]\{0,\}. scripts/conf-2C4C.sh/. scripts/conf-2C4C.sh/g" scripts/test-all.sh
# sed -i '' "s/[\#]\{0,\}. scripts\/conf-noise.sh/\#. scripts\/conf-noise.sh/g" scripts/test-all.sh
# sed -i '' "s/[\#]\{0,\}. scripts\/conf-2C4C.sh/\#. scripts\/conf-2C4C.sh/g" scripts/test-all.sh
# sed -i '' "s/[\#]\{0,\}. scripts\/conf-noise.sh/. scripts\/conf-noise.sh/g" scripts/test-all.sh

# test 4C8C
# sed -i '' "s/CPU_REQUEST=[0-9]\{1,\}/CPU_REQUEST=4/g" scripts/test-all.sh
# sed -i '' "s/CPU_LIMIT=[0-9]\{1,\}/CPU_LIMIT=8/g" scripts/test-all.sh
sed -i '' "s/CPU_REQUEST=\(.*\)/CPU_REQUEST=4/g" scripts/conf-2C4C.sh
sed -i '' "s/CPU_LIMIT=\(.*\)/CPU_LIMIT=8/g" scripts/conf-2C4C.sh

for ((i=0;i<3;i++))
do
    for period in 10 12 13 14 15 25 50 100 200 400 800 1600
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

    for period in 200 400 800 1600
    do
        echo "======+++++    start small $period $i    +++++========="
        while :
        do
            timeout --foreground "$((20*60))" bash scripts/test-all.sh $ENV small $period "$LOG_PAR_DIR"
            if [ $? -eq 0 ]; then
                break
            fi
        done
        echo "======+++++    end   small $period $i    +++++========="
    done
done