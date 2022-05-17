#!/bin/bash

./pub --wsAddr localhost:12001 \
--wsPath /echo \
--peerId peerIdPub \
--roomId roomId2005 \
--appId 3 \
--videoFileName /Downloads/countdown.h264 \
--streamLoop 3 \
--h264FrameDuration 33 \
--rtpOutboundMTU 1200 \
--clockRate 90000 \
--ssrc 56785678 \
--payloadType 125 \
--streamId 123123
