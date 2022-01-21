#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: deploy-logger.sh SERVER PORT"
  exit 1
fi

SERVER=$1
PORT=$2

GOOS=linux GOARCH=amd64 go build -o bin/logserverd-linux-amd64 cmd/logserverd/logserverd.go

rsync bin/logserverd-linux-amd64 $SERVER:projects/
ssh $SERVER "
tmux new -s log -d ; \
tmux send-keys -t log:0.0 C-c ; \
sleep 1;
mkdir -p /var/log/hongbao
tmux send-keys -t log:0.0 '~/projects/logserverd-linux-amd64 -addr 0.0.0.0:$PORT -f /var/log/hongbao' C-m;
"

