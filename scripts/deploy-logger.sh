#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: deploy-logger.sh SERVER PORT"
  exit 1
fi

SERVER=$1
PORT=$2

rsync bin/logserverd-linux-amd64 $SERVER:
ssh $SERVER "tmux new -s log -d"
ssh $SERVER "tmux send-keys -t log:0.0 '/root/logserverd-linux-amd64 -addr 0.0.0.0:$PORT -f /var/log/hongbao' C-m"

