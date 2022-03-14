#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 3 ]; then
  echo "usage: deploy-msd.sh SERVER ENV LOC"
  exit 1
fi

SERVER=$1
ENV=$2
LOC=$3

PRO_DIR='$HOME/projects/hongbao-ms'

rsync -aP ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore --exclude=data

ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-msd.sh $ENV $LOC"
