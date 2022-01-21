#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: deploy-msd.sh SERVER CONFIG"
  exit 1
fi

SERVER=$1
CONFIG=$2

PRO_DIR='$HOME/projects/hongbao-ms'

rsync -aP ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore

ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-msd-cfg.sh $CONFIG"
