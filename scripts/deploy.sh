#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: deploy.sh SERVER ENV LOC"
  exit 1
fi

SERVER=$1
ENV=$2
LOC=$3

SESSION_NAME="ms"
PRO_DIR="~/projects/hongbao-ms"

bash scripts/build.sh
rsync -aP scripts $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash $PRO_DIR/scripts/docker-run.sh $ENV $LOC' C-m"
