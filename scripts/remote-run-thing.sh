#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: remote-run-thing.sh SERVER ENV LOC"
  exit 1
fi

SERVER=$1
ENV=$2
LOC=$3

SESSION_NAME="thing"
PRO_DIR="~/projects/hongbao-ms"

rsync -aP ./* $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash scripts/run-thing.sh $ENV $LOC' C-m"