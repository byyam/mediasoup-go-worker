#!/bin/bash
set -x


# Set environment variables
#export RTC_STATIC_PORT=50000
#export PIPE_PORT=50001


exec ./server --rtcListenIp 0.0.0.0 \
--logLevel trace \
--rtcStaticPort ${RTC_STATIC_PORT} \
--pipePort ${PIPE_PORT} \
--prometheusPort 15000
