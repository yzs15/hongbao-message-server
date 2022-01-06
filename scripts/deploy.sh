#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: deploy.sh SERVER ENV"
  exit 1
fi

SERVER=$1
ENV=$2

SESSION_NAME="ms"
PRO_DIR="~/projects/hongbao-ms"

bash scripts/build.sh
rsync -aP scripts $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash $PRO_DIR/scripts/docker-run.sh wang $ENV' C-m"

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 'bash $PRO_DIR/scripts/docker-run.sh thing $ENV' C-m"
