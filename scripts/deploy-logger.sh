

SERVER=$1

rsync bin/logserverd-linux-amd64 $SERVER:
ssh $SERVER "tmux new -s log -d"
ssh $SERVER "tmux send-keys -t log:0.0 '/root/logserverd-linux-amd64 -addr 0.0.0.0:5555 -f /var/log/hongbao' C-m"

