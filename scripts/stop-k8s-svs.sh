#!/bin/bash
cd $(dirname "$0")
cd ..
. scripts/utils.sh

if [ $# -lt 1 ]; then
  echo "usage: stop-ks8-svs.sh SERVER SUDO"
  exit 1
fi

SERVER=$1
YML_NAME=$2
SUDO=$3

SESSION_NAME=k8s

PRO_DIR='projects/hongbao-ms'
CMD="bash $PRO_DIR/scripts/stop-k8s-svs-local.sh $YML_NAME $SUDO"
ensure_ok ssh $SERVER "'""$CMD""'"
