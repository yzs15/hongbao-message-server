#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: remote-run-thing.sh SERVER ENV"
  exit 1
fi

SERVER=$1
ENV=$2

SESSION_NAME="thing"
PRO_DIR="~/projects/hongbao-ms"

rsync -aP ./* $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash scripts/run-thing.sh $ENV' C-m"