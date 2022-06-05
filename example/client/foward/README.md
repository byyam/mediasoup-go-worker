

## Forward

---

Forward demonstrates how to use mediasoup-go-worker to forward a media stream of RTP, and play with ffplay.

### Instructions

---

#### Install ffplay and prepare a video file of H264 format

``` shell
ffmpeg -re -stream_loop -1 -i countdown.mp4 -vcodec copy -f rtp rtp://127.0.0.1:12002 > h264.sdp
```

#### Run ffplay with sdp file

``` shell
ffplay -i h264.sdp -protocol_whitelist file,udp,rtp,rtcp
```

#### Run the client

``` shell
go run main.go
```


