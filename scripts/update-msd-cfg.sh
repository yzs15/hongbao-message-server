#!/bin/bash
cd $(dirname "$0")
cd ..
. scripts/utils.sh

if [ $# -lt 2 ]; then
  echo "usage: deploy-msd.sh SERVER CONFIG [NO_PULL]"
  exit 1
fi

SERVER=$1
CONFIG=$2
NO_PULL=$3

PRO_DIR='projects/hongbao-ms'

ensure_ok rsync -a ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore --exclude=data

ensure_ok ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-msd-cfg.sh $CONFIG $NO_PULL"
