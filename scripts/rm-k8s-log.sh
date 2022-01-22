#!/bin/bash
cd $(dirname "$0")
cd ..

if [ $# -lt 1 ]; then
  echo "usage: deploy-logger.sh SERVER"
  exit 1
fi

SERVER=$1

ssh $SERVER "
setopt rmstarsilent
rm -rf /var/log/hongbao/*
"