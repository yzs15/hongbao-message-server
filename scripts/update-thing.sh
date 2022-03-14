#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

if [ $# -lt 4 ]; then
  echo "usage: update-thing.sh SERVER CONFIG NODE NUM NO_PULL"
  exit 1
fi

SERVER=$1
CONFIG=$2
NODE=$3
NUM=$4
NO_PULL=$5

PRO_DIR='$HOME/projects/hongbao-ms'

rsync -aP ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore --exclude=bin --exclude=data

ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-thing.sh $CONFIG $NODE $NUM $NO_PULL"
