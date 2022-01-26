SESSION_NAME=k8s

tmux kill-session -t $SESSION_NAME

tmux new -s $SESSION_NAME -d

tmux split-window -t $SESSION_NAME:0 -h
tmux send-keys -t $SESSION_NAME:0.0 "ssh-tmux.sh lab9 k8s" C-m
tmux send-keys -t $SESSION_NAME:0.1 "ssh-tmux.sh hbnj4 k8s" C-m

tmux a -t $SESSION_NAME