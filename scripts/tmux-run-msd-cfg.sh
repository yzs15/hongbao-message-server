#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: tmux-run-msd-cfg.sh ENV LOC"
  exit 1
fi

CONFIG=$1

SESSION_NAME=msd
PRO_DIR='$HOME/projects/hongbao-ms'

tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? != 0 ]; then
  tmux new-session -s $SESSION_NAME -d
  tmux split-window -t $SESSION_NAME:0
fi

tmux send-keys -t $SESSION_NAME:0.0 C-c C-m
sleep 3
tmux send-keys -t $SESSION_NAME:0.0 C-c C-m
sleep 3
tmux send-keys -t $SESSION_NAME:0.0 "bash $PRO_DIR/scripts/docker-run-msd-cfg.sh $CONFIG" C-m

mkdir -p /var/log/hongbao
#tmux send-keys -t $SESSION_NAME:0.1 C-c C-m
#sleep 1
#tmux send-keys -t $SESSION_NAME:0.1 "$PRO_DIR/bin/logserverd-linux-amd64 -addr 0.0.0.0:8083 -f $PRO_DIR/msd.log" C-m
