# mediasoup-go-worker
Pure Go implementation of [Mediasoup worker](https://github.com/versatica/mediasoup).


<p align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
<a href="Language"><img src="https://img.shields.io/badge/language-golang-green.svg"></a>
 <a href="https://goreportcard.com/report/github.com/byyam/mediasoup-go-worker"><img src="https://goreportcard.com/badge/github.com/byyam/mediasoup-go-worker" alt="Go Report Card"></a>
</p>


### Features
#### Mediasoup worker Protocol

|  version  | stream protocol type |         support         |
|:---------:|:--------------------:|:-----------------------:|
|  < 3.9.0  |      netstring       | :ballot_box_with_check: |
| \>= 3.9.0 |        native        | :ballot_box_with_check: |


```shell
goos: darwin
goarch: amd64
pkg: github.com/byyam/mediasoup-go-worker/pkg/netparser
cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
BenchmarkNetNative_WriteBuffer-12       39985699                29.73 ns/op            4 B/op          1 allocs/op
BenchmarkNetNative_ReadBuffer-12          269317              4487 ns/op              12 B/op          3 allocs/op
BenchmarkNetStrings_WriteBuffer-12       7167288               152.0 ns/op           355 B/op          3 allocs/op
BenchmarkNetStrings_ReadBuffer-12       13219642                88.37 ns/op            4 B/op          1 allocs/op
``

#### Mediasoup worker API

##### worker API

| method                  |         support         |
|:------------------------|:-----------------------:|
| worker.close            | :ballot_box_with_check: |
| worker.dump             | :ballot_box_with_check: |
| worker.getResourceUsage | :ballot_box_with_check: |
| worker.updateSettings   |           WIP           |
| worker.createRouter     | :ballot_box_with_check: |


##### router API

| method                             |         support         |
|:-----------------------------------|:-----------------------:|
| router.close                       | :ballot_box_with_check: |
| router.dump                        | :ballot_box_with_check: |
| router.createWebRtcTransport       | :ballot_box_with_check: |
| router.createPlainTransport        |           WIP           |
| router.createPipeTransport         |           WIP           |
| router.createDirectTransport       |           WIP           |
| router.createActiveSpeakerObserver |           WIP           |
| router.createAudioLevelObserver    |           WIP           |


##### transport API

| method                          |         support         |
|:--------------------------------|:-----------------------:|
| transport.close                 | :ballot_box_with_check: |
| transport.dump                  | :ballot_box_with_check: |
| transport.getStats              | :ballot_box_with_check: |
| transport.connect               | :ballot_box_with_check: |
| transport.setMaxIncomingBitrate |           WIP           |
| transport.setMaxOutgoingBitrate |           WIP           |
| transport.produce               | :ballot_box_with_check: |
| transport.consume               | :ballot_box_with_check: |
| transport.produceData           |           WIP           |
| transport.consumeData           |           WIP           |


##### producer API

| method                    |         support         |
|:--------------------------|:-----------------------:|
| producer.close            | :ballot_box_with_check: |
| producer.dump             | :ballot_box_with_check: |
| producer.getStats         | :ballot_box_with_check: |
| producer.pause            |           WIP           |
| producer.resume           |           WIP           |
| producer.enableTraceEvent |           WIP           |


##### consumer API

| method                      |         support         |
|:----------------------------|:-----------------------:|
| consumer.close              | :ballot_box_with_check: |
| consumer.dump               | :ballot_box_with_check: |
| consumer.getStats           | :ballot_box_with_check: |
| consumer.pause              |           WIP           |
| consumer.resume             |           WIP           |
| consumer.setPreferredLayers |           WIP           |
| consumer.setPriority        |           WIP           |
| consumer.requestKeyFrame    |           WIP           |
| consumer.enableTraceEvent   |           WIP           |


##### dataProducer API

| method                | support |
|:----------------------|:-------:|
| dataProducer.close    |   WIP   |
| dataProducer.dump     |   WIP   |
| dataProducer.getStats |   WIP   |


##### dataConsumer API

| method                                     | support |
|:-------------------------------------------|:-------:|
| dataConsumer.close                         |   WIP   |
| dataConsumer.dump                          |   WIP   |
| dataConsumer.getStats                      |   WIP   |
| dataConsumer.getBufferedAmount             |   WIP   |
| dataConsumer.setBufferedAmountLowThreshold |   WIP   |


##### rtpObserver API

| method                             | support |
|:-----------------------------------|:-------:|
| rtpObserver.close                  |   WIP   |
| rtpObserver.pause                  |   WIP   |
| rtpObserver.resume                 |   WIP   |
| rtpObserver.addProducer            |   WIP   |
| rtpObserver.removeProducer         |   WIP   |


#### Codec


### Usage

#### mediasoup-worker

build out executable file.

``` shell
$ cd cmd/mediasoup-worker
$ go build
```

Replace the worker Binary in [mediasoup](https://github.com/versatica/mediasoup) project.


### References

[mediasoup](https://github.com/versatica/mediasoup)

[mediasoup-go](https://pkg.go.dev/github.com/jiyeyuran/mediasoup-go)
