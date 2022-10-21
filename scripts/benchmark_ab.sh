#!/bin/bash

. ./scripts/env.sh
# env for collecting data
REPORT_PATH="output/${REPORT}.log"
DATA_PATH="output/${REPORT}.csv"
LATENCY_PATH="output/latency_${REPORT}"

n=3000
body=(1 2 3 4)
concurrent=(10)
header=(1)
server_name="ticketing_server"
ports=(8001)
serverIP="http://127.0.0.1"

. ./scripts/build_all.sh

# make folder to store all latency data
if [ ! -d ${LATENCY_PATH} ]; then
  mkdir ${LATENCY_PATH}
fi

# generate request in json for ab
for b in ${body[@]}; do
    python3 ./scripts/ab/generate_request.py ${b}
done

# benchmark
for b in ${body[@]}; do
  for h in ${header[@]}; do
    for c in ${concurrent[@]}; do
        addr="${serverIP}:${ports[i]}"

        # order_server start
        nohup $taskset_more ./output/bin/order_server >>output/log/oreder_nohup.log 2>&1 &
        sleep 1
        echo "order_server running"

        # ticketing_server start
        nohup $taskset_less ./output/bin/${server_name}>>output/log/ticketing-nohup.log 2>&1 &
        sleep 1
        echo "server ${server_name}="ticketing_server"running with $taskset_less"

        # run ab
        echo "Benchmark_Config" >> ${REPORT_PATH}
        echo "${server_name},${c},${b}" >> ${REPORT_PATH}
        latency_file="${LATENCY_PATH}/${server_name}_${c}_${b}.csv"
        $taskset_more ab -k -e ${latency_file} -d -S -q -n ${n} -c ${c}   ${addr}/ping | $tee_cmd

        # stop ticketing_server
        kill -9 $(lsof -t -i:8001)

        # stop order_server
        ps -ef | grep order_server | grep -v grep | awk '{print $2}' | xargs kill -9
        sleep 1
    done
  done
done

# parse data and generate output.csv
python3 ./scripts/ab/parse_data.py ${REPORT_PATH} ${LATENCY_PATH} ${DATA_PATH}

