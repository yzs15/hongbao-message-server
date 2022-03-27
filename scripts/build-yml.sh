YAML_FILEPATH=configs/k8s/numrecd.yml
echo "TOTAL yml " $TOTAL
sed "s/APP_NAME/numrecd-$TOTAL-$((RANDOM))/g" configs/k8s/numrecd-detail.yaml.template > $YAML_FILEPATH
# sed "s/APP_NAME/numrecd-$TOTAL-$((RANDOM))/g" configs/k8s/numrecd-each-node-detail.yaml.template > $YAML_FILEPATH
sed -i '' "s/CPU_REQUEST/$CPU_REQUEST/g" $YAML_FILEPATH
sed -i '' "s/CPU_LIMIT/$CPU_LIMIT/g" $YAML_FILEPATH
sed -i '' "s/CPU_SCALE_THRESHOLD/$CPU_SCALE_THRESHOLD/g" $YAML_FILEPATH
sed -i '' "s/CPU_SCALE_UP_LIMIT/$CPU_SCALE_UP_LIMIT/g" $YAML_FILEPATH
sed -i '' "s/INIT_REPLICA_NUM/$INIT_REPLICA_NUM/g" $YAML_FILEPATH
sed -i '' "s/MAX_REPLICA_NUM/$MAX_REPLICA_NUM/g" $YAML_FILEPATH