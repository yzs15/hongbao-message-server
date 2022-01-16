#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

if [ $# -lt 4 ]; then
  echo "usage: deploy-thing.sh SERVER CONFIG NODE NUM"
  exit 1
fi

SERVER=$1
CONFIG=$2
NODE=$3
NUM=$4

PRO_DIR='$HOME/projects/hongbao-ms'

bash scripts/build.sh

ssh $SERVER "mkdir -p $PRO_DIR"
rsync -aP ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore

ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-thing.sh $CONFIG $NODE $NUM"
