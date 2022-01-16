#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: tmux-run-msd.sh ENV LOC"
  exit 1
fi

ENV=$1
LOC=$2

SESSION_NAME=msd
PRO_DIR='~/projects/hongbao-ms'

tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? != 0 ]; then
  tmux new-session -s $SESSION_NAME -d
fi

tmux send-keys -t $SESSION_NAME:0.0 C-c
tmux send-keys -t $SESSION_NAME:0.0 "bash $PRO_DIR/scripts/docker-run-msd.sh $ENV $LOC" C-m