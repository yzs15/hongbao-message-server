#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

SERVER="kbnj1"
PRO_DIR="~/projects/hongbao-ms"

SESSION_NAME="ms"

bash scripts/build.sh
rsync -aP scripts $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash $PRO_DIR/scripts/docker-run.sh wang' C-m"

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 'bash $PRO_DIR/scripts/docker-run.sh thing' C-m"
