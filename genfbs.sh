#!/bin/bash

rm -rf fbs/FBS
flatc -g --go-module-name github.com/byyam/mediasoup-go-worker/fbs -o fbs fbs/*.fbs

# import cycle not allowed
# flatc -g --gen-object-api --go-module-name github.com/byyam/mediasoup-go-worker/fbs -o fbs fbs/*.fbs

