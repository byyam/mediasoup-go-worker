package sdk

import (
	"log"
	"net"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/byyam/mediasoup-go-worker/h264"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/pkg/errors"
)

type PublishOpt struct {
	StreamId    uint64
	SSRC        uint32
	PayloadType uint8
	DtlsClient  bool // default: server
	ClockRate   uint32
	MimeType    string
}

func (c *Client) Publish(opt PublishOpt) (net.Conn, error) {
	mediasoupFPs, err := c.prepareDtls(opt.DtlsClient)
	if err != nil {
		return nil, err
	}

	rtpCodecParameters := []*mediasoupdata.RtpCodecParameters{
		{
			MimeType:    opt.MimeType,
			PayloadType: opt.PayloadType,
			ClockRate:   int(opt.ClockRate),
			Channels:    0,
			Parameters: mediasoupdata.RtpCodecSpecificParameters{
				RtpParameter: h264.RtpParameter{
					PacketizationMode:     1,
					LevelAsymmetryAllowed: 1,
				},
				ProfileId:           "",
				Apt:                 0,
				SpropStereo:         0,
				Useinbandfec:        0,
				Usedtx:              0,
				Maxplaybackrate:     0,
				XGoogleMinBitrate:   0,
				XGoogleMaxBitrate:   0,
				XGoogleStartBitrate: 0,
				ChannelMapping:      "",
				NumStreams:          0,
				CoupledStreams:      0,
			},
			RtcpFeedback: []mediasoupdata.RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
	}
	headerExtensions := []mediasoupdata.RtpHeaderExtensionParameters{
		{
			Uri:     "urn:ietf:params:rtp-hdrext:sdes:mid",
			Id:      4,
			Encrypt: false,
		},
	}
	encodings := []*mediasoupdata.RtpEncodingParameters{
		{
			Ssrc: opt.SSRC,
		},
	}
	rtpParameters := mediasoupdata.RtpParameters{
		Mid:              "0",
		Codecs:           rtpCodecParameters,
		HeaderExtensions: headerExtensions,
		Encodings:        encodings,
		Rtcp:             mediasoupdata.RtcpParameters{},
	}
	webRtcTransportOffer := isignal.WebRtcTransportOffer{
		ForceTcp: false,
		DtlsParameters: mediasoupdata.DtlsParameters{
			Role:         mediasoupdata.DtlsRole(c.dtlsRole),
			Fingerprints: mediasoupFPs,
		},
	}
	req := isignal.PublishRequest{
		StreamId:    opt.StreamId,
		TransportId: "",
		Offer:       webRtcTransportOffer,
		PublishOffer: isignal.PublishOffer{
			Kind:          mediasoupdata.MediaKind_Video,
			RtpParameters: rtpParameters,
			AppData:       nil,
		},
	}
	rsp, err := wsconn.NewWsClient(c.wsOpt).Publish(req)
	if err != nil {
		log.Println("publish response error:", err)
		return nil, err
	}
	c.webrtcTransportAnswer = &rsp.Answer

	conn, err := c.Conn()
	if err != nil {
		return nil, errors.WithMessage(err, "connect sfu error")
	}
	log.Println("handle publish completed")
	return conn, nil
}
