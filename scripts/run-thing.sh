#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: run-thing.sh ENV"
  exit 1
fi

ENV=$1

prefix=$(date +"%Y-%m-%d %H:")

cur_min=$(date +"%-M")
nxt_min=$((cur_min + 1))
if [ $nxt_min -lt 10 ]; then
  nxt_min="0"$nxt_min
fi

start=$prefix$nxt_min:00

if [ "net" = $ENV ]; then
  THING_MS_IP=172.16.32.12
else
  THING_MS_IP=192.168.143.1
fi

go run cmd/thingcli/thingcli.go -start "$start" \
    -duration 30m \
    -ws   $THING_MS_IP:5544 \
    -zmq  tcp://$THING_MS_IP:5543