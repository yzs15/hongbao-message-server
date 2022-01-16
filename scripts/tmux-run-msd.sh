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
PRO_DIR='$HOME/projects/hongbao-ms'

tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? != 0 ]; then
  tmux new-session -s $SESSION_NAME -d
  tmux split-window -t $SESSION_NAME:0
fi

tmux send-keys -t $SESSION_NAME:0.0 C-c C-m
tmux send-keys -t $SESSION_NAME:0.0 "bash $PRO_DIR/scripts/docker-run-msd.sh $ENV $LOC" C-m

mkdir -p /var/log/hongbao
tmux send-keys -t $SESSION_NAME:0.1 "touch $PRO_DIR/msd.log" C-m
tmux send-keys -t $SESSION_NAME:0.1 "$PRO_DIR/bin/logserverd-linux-amd64 -addr 0.0.0.0:8083 -f /var/log/hongbao -f $PRO_DIR/msd.log" C-m