#!/bin/bash

. ./scripts/env.sh

ports=8001

. ./scripts/build_all.sh
n=3000

body=(1)
concurrent=(100)
serverIP="http://127.0.0.1"

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
      addr="${serverIP}:${ports}"
      echo $taskset_more ./output/bin/ticketing_server

		
      # order_server start
      nohup $taskset_more ./output/bin/order_server >>output/log/oreder_nohup.log 2>&1 &
      sleep 1
      echo "order_server running"

      # ticketing_server start
      nohup $taskset_more ./output/bin/ticketing_server >>output/log/ticketing_nohup.log 2>&1 &
      sleep 1
      echo "ticketing_server running"

      # run client
      echo "client running"
      $taskset_less ./output/bin/client -addr="$addr"/ping -b=$b -c=$c -n=$n -s=ticketing_server | $tee_cmd

      # stop ticketing_server
     kill -9 $(lsof -t -i:8001)
     # stop order_server
     ps -ef | grep order_server | grep -v grep | awk '{print $2}' | xargs kill -9
     # stop client
     ps -ef | grep client/client | grep -v grep | awk '{print $2}' | xargs kill -9	
      sleep 1
  done
done
