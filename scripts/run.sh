#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: run.sh NAME"
  exit 1
fi

NAME=$1

if [ "wang" = $NAME ]; then
  go run cmd/msd/msd.go -wang \
      -net

elif [ "thing" = $NAME ]; then
  go run cmd/msd/msd.go -thing \
      -ws   0.0.0.0:5544 \
      -zmq  tcp://0.0.0.0:5543 \
      -log  0.0.0.0:5542 \
      -net \
      -wend tcp://127.0.0.1:5553 \
      -kend 172.16.32.12:32101 \
      -kend 172.16.32.13:32101 \
      -kend 172.16.32.14:32101 \
      -kend 172.16.32.15:32101

fi