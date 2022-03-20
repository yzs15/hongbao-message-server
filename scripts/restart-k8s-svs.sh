#!/bin/bash
cd $(dirname "$0")
cd ..
. scripts/utils.sh

if [ $# -lt 1 ]; then
  echo "usage: deploy-logger.sh SERVER SUDO"
  exit 1
fi

SERVER=$1
YML_NAME=$2
SUDO=$3

SESSION_NAME=k8s

PRO_DIR='projects/hongbao-ms'
ensure_ok ssh $SERVER "'""bash $PRO_DIR/scripts/restart-k8s-svs-local.sh $YML_NAME $SUDO""'"
# tmux send-keys -t $SESSION_NAME:0.0 '$SUDO kubectl apply -f ~/projects/k8s/hongbaod.yml' C-m