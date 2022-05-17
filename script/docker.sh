#!/bin/bash
set -x

if [ $# -lt 1 ];then
        echo usage "$0 [build|run]"
        exit 1
fi

ACTION=$1
RTC_STATIC_PORT=50000
PIPE_PORT=50001

if  [[ "${ACTION}" == "build" ]];then
    docker build -t mediasoupworker -f deploy/Dockerfile . --build-arg RTC_STATIC_PORT=${RTC_STATIC_PORT} --build-arg PIPE_PORT=${PIPE_PORT}
elif [[ "${ACTION}" == "run" ]];then
    docker run -d -p 127.0.0.1:$PIPE_PORT:$PIPE_PORT/udp -p 127.0.0.1:$RTC_STATIC_PORT:$RTC_STATIC_PORT/udp -p 127.0.0.1:12002:12002/tcp mediasoupworker
else
    echo usage "$0 [build|run]"
    exit 1
fi



