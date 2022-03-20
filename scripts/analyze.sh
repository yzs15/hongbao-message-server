#!/bin/bash
cd $(dirname "$0")
cd ..

cd ../hongbao-log
source ./venv/bin/activate

DIR_PATH=$1

LOG_DIR=$DIR_PATH/$(ls -l $DIR_PATH | grep 2022 | awk '{print $9}')
K8S_CPU=$DIR_PATH"/k8s-cpu"
TS_CPU=$DIR_PATH"/ts-cpu"

if [ -d "$K8S_CPU" ]; then
    echo "cal k8s cpu usage ......"
    python3 ./src/k8s_cpu_utilization.py $K8S_CPU
fi

if [ -d "$TS_CPU" ]; then
    echo "cal ts cpu alloc ......"
    python3 ./src/ts_cpu_alloc.py $LOG_DIR
    echo "cal ts cpu usage 1ms ......"
    python3 ./src/ts_cpu_usage.py $TS_CPU 1
    echo "cal ts cpu usage 100ms ......"
    python3 ./src/ts_cpu_usage.py $TS_CPU 100
fi

