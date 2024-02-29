#!/bin/bash

rm -rf fbs/FBS

# cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release
# make && sudo make install
flatc -g --gen-object-api --go-module-name github.com/byyam/mediasoup-go-worker/fbs -o fbs fbs/*.fbs

