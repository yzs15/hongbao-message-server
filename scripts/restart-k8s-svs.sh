#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: deploy-logger.sh SERVER SUDO"
  exit 1
fi

SERVER=$1
SUDO=$2

SESSION_NAME=k8s

ssh $SERVER "
tmux send-keys -t $SESSION_NAME:0.0 '$SUDO kubectl delete -f ~/projects/k8s/numrecd.yml' C-m
tmux send-keys -t $SESSION_NAME:0.0 '$SUDO kubectl delete -f ~/projects/k8s/hongbaod.yml' C-m
tmux send-keys -t $SESSION_NAME:0.0 '$SUDO rm -rf /var/log/hongbao/*' C-m
tmux send-keys -t $SESSION_NAME:0.0 '$SUDO kubectl apply -f ~/projects/k8s/numrecd.yml' C-m
tmux send-keys -t $SESSION_NAME:0.0 '$SUDO kubectl apply -f ~/projects/k8s/hongbaod.yml' C-m
"