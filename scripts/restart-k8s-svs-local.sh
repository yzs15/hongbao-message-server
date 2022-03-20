SESSION_NAME=k8s

YML_NAME=$1
SUDO=$2

tmux send-keys -t $SESSION_NAME:0.0 C-c C-m ;
sleep 1
tmux send-keys -t $SESSION_NAME:0.0 "$SUDO kubectl delete -f ~/projects/k8s/$YML_NAME.yml" C-m ;
sleep 1
tmux send-keys -t $SESSION_NAME:0.0 "$SUDO rm -rf /var/log/hongbao/*" C-m ;
sleep 1
tmux send-keys -t $SESSION_NAME:0.0 "$SUDO kubectl apply -f ~/projects/k8s/$YML_NAME.yml" C-m ;
sleep 1
tmux send-keys -t $SESSION_NAME:0.0 C-c C-m ;