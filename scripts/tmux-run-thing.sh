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
#if [ $? = 0 ]; then
#  PANE_NUM=$(tmux list-panes -t $SESSION_NAME | wc -l)
#  if [ $PANE_NUM -lt $NUM ]; then
    tmux kill-session -t $SESSION_NAME
    tmux new-session -s $SESSION_NAME -d
    for (( i=1; i<=($NUM+3)/4; i++ ))
    do
      tmux new-window -t $SESSION_NAME:$i
      tmux split-window -t $SESSION_NAME:$i
      tmux split-window -t $SESSION_NAME:$i
      tmux split-window -t $SESSION_NAME:$i
      tmux select-layout -t $SESSION_NAME:$i tiled
    done
    tmux kill-window -t $SESSION_NAME:0
#  fi
#fi

idx=0
for (( wi=1; wi<=($NUM+3)/4; wi++ ))
do
  for (( pi=0; pi<4; pi++ ))
  do
    tmux send-keys -t $SESSION_NAME:$wi.$pi C-c C-m
    tmux send-keys -t $SESSION_NAME:$wi.$pi "bash $PRO_DIR/scripts/docker-run-thing.sh $CONFIG $NODE $idx" C-m

    (( idx++ ))
    if [[ $idx -eq $NUM ]]; then
      break
    fi
  done

  if [[ $idx -eq $NUM ]]; then
    break
  fi
done