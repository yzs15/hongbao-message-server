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

GOOS=linux GOARCH=amd64 go build -o bin/logserverd-linux-amd64 cmd/logserverd/logserverd.go
bash scripts/build.sh

ssh $SERVER "mkdir -p $PRO_DIR"
rsync -aP ./* $SERVER:$PRO_DIR/  --exclude-from=.gitignore

ssh $SERVER "bash $PRO_DIR/scripts/tmux-run-msd-cfg.sh $CONFIG"
