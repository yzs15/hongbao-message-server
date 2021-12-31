#!/bin/bash
cd $(dirname "$0")
cd ..

NAME=$1

if [ "wang" = $NAME ]; then
  go run cmd/msd/msd.go --wang

elif [ "thing" = $NAME ]; then
  go run cmd/msd/msd.go --thing \
      --ws   0.0.0.0:5544 \
      --zmq  tcp://0.0.0.0:5543 \
      --log  0.0.0.0:5542 \
      --wend tcp://127.0.0.1:5553 \
      --kend 172.16.32.12:32101 \
      --kend 172.16.32.13:32101 \
      --kend 172.16.32.14:32101 \
      --kend 172.16.32.15:32101

else
  echo "usage: run.sh NAME"
fi