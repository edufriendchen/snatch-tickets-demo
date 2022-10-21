#!/bin/bash

. ./scripts/env.sh
# env for collecting data
REPORT_PATH="output/${REPORT}.log"
DATA_PATH="output/${REPORT}.csv"

t=10
body=(1 2 4 8)
concurrent=(100)
header=(1)
service_name="ticketing_server"
ports=(8001)
serverIP="http://127.0.0.1"

. ./scripts/build_all.sh

# benchmark
for b in ${body[@]}; do
  for h in ${header[@]}; do
    for c in ${concurrent[@]}; do
        addr="${serverIP}:${ports[i]}"

        # server start
        nohup $taskset_less ./output/bin/${service_name} >>output/log/nohup.log 2>&1 &
        sleep 1
        echo "${service_name}  running with $taskset_less"

        # run wrk
        echo "Benchmark_Config" >> ${REPORT_PATH}
        echo "${service_name},${c},${b}" >> ${REPORT_PATH}
        $taskset_more wrk -d${t}s -s  ./scripts/wrk/benchmark.lua -c${c} -t${c} ${addr}/ping -- ${b} | $tee_cmd

        # stop server
        pid=$(ps -ef | grep ${service_name} | grep -v grep | awk '{print $2}')
        disown $pid
        kill -9 $pid
        sleep 1
    done
  done
done

# parse data and generate output.csv
python3 ./scripts/wrk/parse_data.py ${REPORT_PATH} ${DATA_PATH}

# 使用数据生产性能图
python3 ./scripts/reports/render_images.py ${REPORT}





