#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 2 ]; then
  echo "usage: run.sh ENV LOC"
  exit 1
fi

ENV=$1
LOC=$2

if [ "net" = $ENV ]; then
  if [ "bj" = $LOC ]; then
    KENDS="-kend 10.208.104.9"
    NSEND="58.213.121.2:10027"
    ZMQ_OUT="tcp://159.226.41.229:7102"
  else
    KENDS="-kend 172.16.32.13 -kend 172.16.32.14 -kend 172.16.32.15"
    NSEND="172.16.32.13:8080"
    ZMQ_OUT="tcp://58.213.121.2:10025"
  fi
else
  KENDS=""
fi

go run cmd/msd/*.go \
    -ws    0.0.0.0:5554 \
    -zmq   tcp://0.0.0.0:5553 \
    -log   0.0.0.0:5552 \
    -log-path $PWD \
    -nsend $NSEND \
    -zmq-out $ZMQ_OUT \
    -$ENV  $KENDS
