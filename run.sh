#!/bin/bash

MEDIASOUP_VERSION="3.9.0" ./bin/mediasoup-worker-darwin --logLevel=warn --logTags=info --logTags=ice --logTags=dtls --logTags=rtp --logTags=srtp \
--logTags=rtcp --logTags=rtx --logTags=bwe --logTags=score --logTags=simulcast --logTags=svc --logTags=sctp \
--rtcMinPort=40000 --rtcMaxPort=49999 --prometheusPath=/metrics --prometheusPort=15000

