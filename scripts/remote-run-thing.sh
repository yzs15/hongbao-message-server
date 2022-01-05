#!/bin/bash
cd $(dirname "$0")
cd ..

SERVER="kbnj1"
PRO_DIR="~/projects/hongbao-ms"

SESSION_NAME="thing"

rsync -aP ./* $SERVER:$PRO_DIR/

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.0 'bash scripts/run-thing.sh' C-m"

ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 C-c"
ssh $SERVER "tmux send-keys -t $SESSION_NAME:0.1 'date' C-m"
