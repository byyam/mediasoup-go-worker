#!/bin/bash

./pkglogger --rtcListenIp 127.0.0.1 \
--logLevel trace \
--rtcStaticPort 50000 \
--pipePort 50001 \
--prometheusPort 15000

