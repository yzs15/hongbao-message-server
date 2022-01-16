#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: tmux-run-thing.sh CONFIG NODE NUM"
  exit 1
fi

CONFIG=$1
NODE=$2
NUM=$3

SESSION_NAME=thing
PRO_DIR='$HOME/projects/hongbao-ms'

tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? = 0 ]; then
  tmux kill-session -t $SESSION_NAME
fi

tmux new-session -s $SESSION_NAME -d
for (( i=1; i<$NUM; i++ ))
do
  tmux split-window -t $SESSION_NAME:0
done
tmux select-layout -t $SESSION_NAME:0 tiled

for (( i=0; i<$NUM; i++))
do
  tmux send-keys -t $SESSION_NAME:0.$i "bash $PRO_DIR/scripts/docker-run-thing.sh $CONFIG $NODE $i" C-m
done