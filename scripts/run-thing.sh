#!/bin/bash
set -x
cd $(dirname "$0")
cd ..

prefix=$(date +"%Y-%m-%d %H:")

cur_min=$(date +"%-M")
nxt_min=$((cur_min + 1))
if [ $nxt_min -lt 10 ]; then
  nxt_min="0"$nxt_min
fi

start=$prefix$nxt_min:00

go run cmd/thingcli/thingcli.go -start "$start" \
    -duration 30m