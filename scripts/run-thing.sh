#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: run-thing.sh ENV LOC"
  exit 1
fi

ENV=$1
LOC=$2

if [ "net" = $ENV ]; then
  if [ "bj" = $LOC ]; then
    MS_WS_END=159.226.41.229:7101
    MS_ZMQ_END=tcp://159.226.41.229:7102
    MAC_ADDR=02:42:ac:12:00:00
  else
    MS_WS_END=172.16.32.12:8080
    MS_ZMQ_END=tcp://172.16.32.12:5557
    MAC_ADDR=02:42:ac:12:00:01
  fi
else
  MS_WS_END=58.213.121.2:10022
  MS_ZMQ_END=tcp://58.213.121.2:10025
  if [ "bj" = $LOC ]; then
    MAC_ADDR=02:42:ac:12:00:02
  else
    MAC_ADDR=02:42:ac:12:00:03
  fi
fi

prefix=$(date +"%Y-%m-%d %H:")
cur_min=$(date +"%-M")
nxt_min=$((cur_min + 1))
if [ $nxt_min -lt 10 ]; then
  nxt_min="0"$nxt_min
fi
START=$prefix$nxt_min:00

go run cmd/thingcli/thingcli.go -start "$START" \
    -duration 30s \
    -ws   $MS_WS_END \
    -zmq  $MS_ZMQ_END \
    -mac  $MAC_ADDR