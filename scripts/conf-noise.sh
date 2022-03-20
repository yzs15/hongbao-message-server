# K8s 容器 配置
export CPU_REQUEST=100m
export CPU_LIMIT=200m
export CPU_SCALE_THRESHOLD=60
export CPU_SCALE_UP_LIMIT=1000
export INIT_REPLICA_NUM=1
export MAX_REPLICA_NUM=1000

# 物端提交配置
export WITH_NOISE=TRUE
export TOTAL_TIME=40
## 每周期负载发送的数量
export LOAD_NUM_PER=1
## 每周期噪音发送的数量
export NOISE_NUM_PER=1

# 分析器分析开始时间
export TEST_TIME_LENGTH=100